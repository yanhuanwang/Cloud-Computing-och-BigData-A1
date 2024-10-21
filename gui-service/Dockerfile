FROM nginx:latest

# Copy the HTML files to the Nginx HTML directory
COPY src/index.html /usr/share/nginx/html/

# Expose port 80 to allow external access to the Nginx server
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
