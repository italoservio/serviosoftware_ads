# ServioSoftware Ads - Cloaker System

## How It Works

When someone accesses a cloaker URL like:
```
http://localhost:8080/r/MTIzNDU2Nzg5MGFiY2RlZg==
```

The system will:

1. **Extract the encoded ID** from the URL
2. **Get the visitor's real IP address** (handling proxies via X-Forwarded-For, X-Real-IP)
3. **Get the User-Agent** from the request headers
4. **Execute the cloaking logic:**
   - Check if bot → white page
   - Check mobile-only config → white page if not mobile
   - Check IP lookup cache → use cached decision
   - Fetch IP metadata from Netify (if not cached)
   - Check known applications (Facebook, TikTok, etc.) → white page
   - Check shared score < 60 → white page (likely bot/datacenter)
   - Otherwise → **black page** (real offer)
5. **Perform HTTP 302 redirect** to the appropriate URL

The redirect is **fast and public** (no authentication required), perfect for ad campaigns! 🚀

## Architecture

- **MongoDB** - Stores cloakers and IP lookup cache
- **Netify API** - Provides IP metadata (shared score, applications, geolocation)
- **JWT Authentication** - Protects management endpoints
- **Gorilla Mux** - HTTP routing

## Key Features

- ✅ Bot detection via User-Agent analysis
- ✅ Mobile-only filtering
- ✅ IP pattern caching (reduces API calls)
- ✅ Known application detection (Facebook, Google, TikTok, etc.)
- ✅ Shared score analysis (dedicated vs shared IPs)
- ✅ Async IP lookup creation (doesn't slow down redirects)
- ✅ Access count tracking
