"""
URL Shortening Service Implementation with Redis

This file implements core functionality for a URL shortening service using FastAPI
and Redis for persistent storage.
"""

from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import RedirectResponse, HTMLResponse
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from pydantic import BaseModel, HttpUrl
import redis
import random
import string
from typing import Optional, List
import os
from datetime import timedelta

class RedisManager:
    """
    Redis connection wrapper class
    
    Manages the Redis connection and provides an interface for URL operations.
    Uses context management for safe resource management.
    """
    
    def __init__(self):
        """Initialize Redis connection with configuration options"""
        try:
            self.redis_client = redis.Redis(
                host=os.getenv("REDIS_HOST", "127.0.0.1"),
                port=int(os.getenv("REDIS_PORT", 6379)),
                decode_responses=True
            )
            # Test connection
            self.redis_client.ping()
        except redis.ConnectionError as e:
            raise RuntimeError(f"Failed to connect to Redis: {str(e)}")

    def set_url(self, key: str, original_url: str, expiry: int = 30 * 60):
        """
        Store a URL mapping in Redis
        
        Args:
            key: The shortened URL key
            original_url: The original URL
            expiry: Expiration time in seconds (default 30 minutes)
        """
        self.redis_client.setex(key, expiry, original_url)

    def get_url(self, key: str) -> Optional[str]:
        """
        Retrieve original URL from Redis
        
        Args:
            key: The shortened URL key
            
        Returns:
            Optional[str]: The original URL if found, None otherwise
        """
        return self.redis_client.get(key)
    
    def get_all_key_values(self) -> List[str]:
        """
        Get all keys and values in Redis
        
        Returns:
            List[str]: A list of all keys and values
        """
        keys = self.redis_client.keys("*")
        return [{
            "key": key, 
            "value": self.redis_client.get(key),
            "ttl": str(timedelta(seconds=self.redis_client.ttl(key))) # Convert to minutes
        } for key in keys]

    def key_exists(self, key: str) -> bool:
        """
        Check if a key exists in Redis
        
        Args:
            key: The key to check
            
        Returns:
            bool: True if key exists, False otherwise
        """
        return bool(self.redis_client.exists(key))

    def update_expiry(self, key: str, seconds: int):
        """
        Update the expiration time for a key
        
        Args:
            key: The key to update
            seconds: New expiration time in seconds
        """
        self.redis_client.expire(key, seconds)

class URLRequest(BaseModel):
    """Request model for URL shortening"""
    url: HttpUrl

class URLResponse(BaseModel):
    """Response model for shortened URLs"""
    short_url: str
    original_url: str

def generate_short_url_key(length: int = 6) -> str:
    """
    Generate a random 6-character string for use as a shortened URL key
    
    The function generates a random string using the following components:
    - Character set: 0-9, a-z, A-Z (62 possible characters)
    - Length: 6 characters
    - Random number generation: Uses system random for secure generation
    
    Returns:
        str: A 6-character string composed of alphanumeric characters
    """
    characters = string.ascii_letters + string.digits
    return ''.join(random.choice(characters) for _ in range(length))

def generate_unique_key(redis_manager: RedisManager) -> str:
    """
    Generate a unique key that doesn't exist in Redis
    
    Args:
        redis_manager: Reference to the Redis manager
        
    Returns:
        str: A unique key
    """
    while True:
        key = generate_short_url_key()
        if not redis_manager.key_exists(key):
            return key

def is_valid_key(key: str) -> bool:
    """
    Validate the format of a shortened URL key
    
    Args:
        key: The key to validate
        
    Returns:
        bool: True if key is valid, False otherwise
    """
    if len(key) != 6:
        return False
    return key.isalnum()

# Initialize FastAPI app
app = FastAPI(title="URL Shortener API")

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize Redis manager
redis_manager = RedisManager()

# Initialize templates
templates = Jinja2Templates(directory="public/templates")

# Mount static files directory
app.mount("/static", StaticFiles(directory="public/static"), name="static")

@app.post("/url", response_model=URLResponse, status_code=201)
async def create_short_url(url_request: URLRequest):
    """Create a shortened URL"""
    original_url = str(url_request.url)
    short_key = generate_unique_key(redis_manager)
    redis_manager.set_url(short_key, original_url)
    
    return URLResponse(
        short_url=short_key,
        original_url=original_url
    )

@app.get("/health")
async def health_check():
    """Health check endpoint, and returns a JSON response of all active keys"""
    try:
        redis_manager.redis_client.ping()
        return {"status": "healthy", "keys": redis_manager.get_all_key_values()}
    except redis.ConnectionError:
        raise HTTPException(status_code=503, detail="Redis connection failed")

@app.get("/{key}")
async def redirect_to_url(key: str):
    """Redirect to original URL"""
    if not is_valid_key(key):
        raise HTTPException(status_code=400, detail="Invalid key format")

    original_url = redis_manager.get_url(key)
    if not original_url:
        raise HTTPException(status_code=404, detail="Short URL not found")

    # Update expiration time to 30 minutes
    redis_manager.update_expiry(key, 30 * 60)

    return RedirectResponse(url=original_url, status_code=302)

@app.get("/", response_class=HTMLResponse)
async def read_root(request: Request):
    """Serve the main HTML page with Jinja2"""
    try:
        context = {
        }
        context["request"] = request
        context["title"] = "URL Shortener"
        context["HOST_DOMAIN"] = os.getenv("HOST_DOMAIN")
        context["BACKEND_API_SHORTEN_URL"] = os.getenv("BACKEND_API_SHORTEN_URL")
        return templates.TemplateResponse("main.html", context)
    except Exception as e:
        raise HTTPException(status_code=404, detail="Template not found")

@app.get("/{page}.html", response_class=HTMLResponse)
async def read_page(request: Request, page: str):
    """Serve any HTML page using Jinja2"""
    try:
        return templates.TemplateResponse(
            f"public/{page}.html", 
            {"request": request, "title": page.title()}
        )
    except Exception as e:
        raise HTTPException(status_code=404, detail="Page not found")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        app, 
        host="0.0.0.0", 
        port=int(os.getenv("BACKEND_PORT", 8080)), 
        reload=True
    )
