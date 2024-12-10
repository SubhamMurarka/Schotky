FROM openresty/openresty

# Install necessary packages for LuaRocks and other dependencies
RUN apt-get update

RUN apt-get install -y luarocks build-essential

# Install Redis client for Lua
RUN luarocks install lua-resty-redis

RUN luarocks install murmurhash3