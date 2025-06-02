#!/bin/bash

# Menyimpan password
PASSWORD="ERPLabIIM-007"

# First ensure the target directory exists and has proper permissions
sshpass -p "$PASSWORD" ssh erplabiim@103.127.134.4 "sudo mkdir -p /var/www/lecsens-new/be-lecsens/asset_management"

# Menyinkronkan direktori lokal ke server menggunakan sshpass dan rsync with sudo
# Gunakan path WSL yang benar dan current directory
sshpass -p "$PASSWORD" rsync -avz --rsync-path="sudo rsync" ./ erplabiim@103.127.134.4:/var/www/lecsens-new/be-lecsens/asset_management