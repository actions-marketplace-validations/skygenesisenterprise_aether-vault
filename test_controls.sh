#!/bin/bash

echo "Testing vault totp list controls..."
echo "The interface should start. Try these keys:"
echo "q - quit"
echo "c - copy code"
echo "a - add entry (will show not implemented message)"
echo "d - delete entry (will show not implemented message)"
echo ""
echo "Starting vault totp list..."
echo "Press Enter to continue..."
read

cd /home/liam/Bureau/enterprise/aether-vault/package/cli
./vault totp list