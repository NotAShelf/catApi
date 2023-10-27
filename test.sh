#!/usr/bin/env bash

# Directory containing the images
image_directory="images"

# Set the starting number
number=0

# Set the maximum number for renaming (change this to your desired limit)
max_number=100

# Loop through the image files in the directory and rename them
for image in "$image_directory"/*; do
	# Check if we've reached the maximum number
	if [ "$number" -gt "$max_number" ]; then
		echo "Reached the maximum number for renaming."
		break
	fi

	# Get the file extension
	extension="${image##*.}"

	# Rename the image with the current number
	new_name="${image_directory}/${number}.${extension}"
	mv "$image" "$new_name"

	# Increment the number
	number=$((number + 1))
done
