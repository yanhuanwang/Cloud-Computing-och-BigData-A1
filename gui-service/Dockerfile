# Step 1: Build Flask service
FROM python:3.9-slim

# Set working directory
WORKDIR /app

# Copy Flask server code
# COPY flask_gui_server.py /app/
# COPY listExpense.html /app/
COPY src/ /app/
COPY swagger-ui/ /app/

# Install dependencies
RUN pip install flask

# Expose Flask port
EXPOSE 5000

# Start Flask server
CMD ["python", "flask_gui_server.py"]
