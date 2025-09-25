#!/bin/bash

#  Update packages
# sudo apt update
# sudo apt install -y --no-install-recommends \
#     uuid-runtime


# Source NVM to make it available (using system installation path)
export NVM_DIR="/usr/local/share/nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion

# make NVM available and install node -lts
nvm install --lts
uv sync
