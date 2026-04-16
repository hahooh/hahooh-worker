# Use the official Go 1.26 image (Debian Bookworm-based)
FROM golang:1.26-bookworm

# 1. Optimized Node.js 22.x (Current LTS in 2026) & Git installation
# NodeSource changed their URL structure recently; 'nodistro' is now the standard.
RUN apt-get update && apt-get install -y curl git gnupg \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list \
    && apt-get update && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# 2. Install Gemini CLI v0.32+ (April 2026 version)
RUN npm install -g @google/gemini-cli

# 3. Secure Workspace Setup
RUN useradd -m agentuser
WORKDIR /app

# Ensure correct permissions for the agent and shared volume
RUN mkdir -p /shared /app/workspace && chown -R agentuser:agentuser /app /shared

# 4. Multi-Stage-like Build in one step
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
# Go 1.26 has better CGO performance; keeping it enabled for Gemini CLI native bindings
RUN CGO_ENABLED=1 go build -buildvcs=false -o /usr/local/bin/agent-runner ./cmd/agent

# Switch to the non-root user for security
USER agentuser

# 5. Persistent config for Gemini CLI
# This ensures the agent's "memory" and auth stay persistent if mounted
ENV GEMINI_CONFIG_HOME=/shared/.gemini

RUN git config --global user.name "Hahooh Agent"
RUN git config --global user.email "agent@hahooh.local"

ENTRYPOINT ["agent-runner"]