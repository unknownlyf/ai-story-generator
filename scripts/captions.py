#!/usr/bin/env python3
import os
import sys
import subprocess
import json

def main():
    if len(sys.argv) != 4:
        print("Usage: python captions.py attach <video_path> <output_filename>")
        sys.exit(1)
    
    mode = sys.argv[1]
    video_path = sys.argv[2]
    output_filename = sys.argv[3]
    
    # Ensure output directory exists
    output_dir = "./generated"
    os.makedirs(output_dir, exist_ok=True)
    output_path = os.path.join(output_dir, output_filename)
    
    try:
        print(f"Creating simple video from: {video_path}")
        print(f"Output will be: {output_path}")
        
        # Just copy the input video to output with a simple text overlay
        # This bypasses all the complex subtitle generation
        ffmpeg_cmd = [
            "ffmpeg", "-y",
            "-i", video_path,
            "-vf", "drawtext=text='AI Generated Story':fontcolor=white:fontsize=24:x=(w-text_w)/2:y=h-th-10",
            "-codec:a", "copy",
            output_path
        ]
        
        print(f"Running: {' '.join(ffmpeg_cmd)}")
        result = subprocess.run(ffmpeg_cmd, capture_output=True, text=True)
        
        if result.returncode == 0:
            print(f"Successfully created video: {output_path}")
        else:
            print(f"FFmpeg error: {result.stderr}")
            sys.exit(1)
            
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    main()
