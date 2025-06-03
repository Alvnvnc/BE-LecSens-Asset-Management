#!/bin/bash

# Menyimpan password
PASSWORD="ERPLabIIM-007"

# First ensure the target directory exists and has proper permissions
sshpass -p "$PASSWORD" ssh erplabiim@103.127.134.4 "sudo mkdir -p /var/www/lecsens-new/be-lecsens"

# Menyinkronkan direktori lokal ke server menggunakan sshpass dan rsync with sudo
sshpass -p "$PASSWORD" rsync -avz --rsync-path="sudo rsync" /home/alvn/Documents/playground/kp/be-lecsens/asset_management erplabiim@103.127.134.4:/var/www/lecsens-new/be-lecsens/

