FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Navigate to the cmd directory where main.go is located
WORKDIR /app/cmd

# Expose the port to be consumed outside
EXPOSE 4000

# Command to run the Go application directly
CMD ["go", "run", "."]
