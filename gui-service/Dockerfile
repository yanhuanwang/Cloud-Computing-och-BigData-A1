FROM nginx:latest

# Copy the HTML files to the Nginx HTML directory
COPY swagger-ui/ /usr/share/nginx/html/swagger-ui

# Add custom Nginx configuration to allow requests from all domains
# COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80 to allow external access to the Nginx server
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
