FROM nginx:latest

# Copy custom Nginx configuration
# COPY nginx.conf /etc/nginx/nginx.conf
COPY listExpense.html /usr/share/nginx/html

# Expose Nginx port
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]