FROM golang:1.26-bookworm

# 1. Install Node.js 22.x, Git, and GPG
RUN apt-get update && apt-get install -y curl git gnupg \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list \
    && apt-get update && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# 2. Install Gemini CLI
RUN npm install -g @google/gemini-cli

# 3. Setup User and Workdir
RUN useradd -m -u 1000 agentuser
WORKDIR /app

# Create shared directories and set initial ownership
RUN mkdir -p /shared /app/workspace && chown -R agentuser:agentuser /app /shared

# 4. Build Strategy: Layers for Caching
# Copy only dependency files first
COPY --chown=agentuser:agentuser go.mod go.sum* ./
RUN go mod download

# Now copy the rest of the source
COPY --chown=agentuser:agentuser . .

# Build the binary
RUN CGO_ENABLED=1 go build -buildvcs=false -o /usr/local/bin/agent-runner ./cmd/agent

# 5. Runtime Environment
USER agentuser
ENV GEMINI_CONFIG_HOME=/shared/.gemini

# --- ADD THIS LINE ---
# This prevents Git from ever hanging on a password prompt
ENV GIT_TERMINAL_PROMPT=0

# Git identity for the agentuser
RUN git config --global user.name "Hahooh Agent" \
    && git config --global user.email "agent@hahooh.local"

ENTRYPOINT ["agent-runner"]