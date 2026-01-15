#!/bin/bash
# Aether Vault Password Launcher
# This script will prompt for password to decrypt the encrypted file

ENCRYPTED_FILE="./test.ava"
OUTPUT_DIR="./test_decrypted"

echo "🔐 Aether Vault - Encrypted File"
echo "📁 File: $ENCRYPTED_FILE"
echo ""

# Prompt for password
echo "🔑 Enter password to decrypt:"
read -s password
echo ""

# Attempt to decrypt using vault CLI
if vault decrypt "$ENCRYPTED_FILE" --passphrase --output "$OUTPUT_DIR" <<< "$password"; then
    echo "✅ Decryption successful!"
    echo "📁 Files are available in: $OUTPUT_DIR"
    
    # Ask if user wants to open the decrypted folder
    read -p "🚀 Open decrypted folder? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if command -v xdg-open > /dev/null 2>&1; then
            xdg-open "$OUTPUT_DIR"
        elif command -v open > /dev/null 2>&1; then
            open "$OUTPUT_DIR"
        elif command -v explorer > /dev/null 2>&1; then
            explorer "$OUTPUT_DIR"
        else
            echo "📂 Decrypted files available at: $OUTPUT_DIR"
        fi
    fi
else
    echo "❌ Decryption failed! Invalid password or corrupted file."
    echo "Please try again with the correct password."
    exit 1
fi
