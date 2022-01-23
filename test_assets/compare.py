import cv2
from skimage import metrics
import os 
import json
import argparse

def write(d):
    print(json.dumps(d))

def compare(path_1, path_2):
    if not os.path.exists(path_1) or not os.path.exists(path_2):
        write({"error": "no file found in the given path"})
        return

    cap1 = cv2.VideoCapture(path_1)
    cap2 = cv2.VideoCapture(path_2)

    i, total_score = 0, 0
    while cap1.isOpened() and cap2.isOpened():
        res1, frame1 = cap1.read()
        res2, frame2 = cap2.read()      

        if (res1 and not res2) or (not res1 and res2):
            print(json.dumps({"error":"videos does not have equal length"}))
            return
        
        if not res1 and not res2:
            break

        if frame1.shape[2] != 3 or frame2.shape[2] != 3:
            write({"error":"videos are not rgb"})
        
        score, diff = metrics.structural_similarity(frame1, frame2, full=True,  channel_axis=2)
        i+=1
        total_score += score
    
    if i == 0:
        write({"error": "video file cannot be opened"})
        return

    write({"avg_sim":total_score/i})


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-v", "--Videos", nargs = "+", help = "Path of the video files to compare")
    args = parser.parse_args()

    if args.Videos is None or len(args.Videos) < 2:
        print("invalid arguments")

    compare(args.Videos[0], args.Videos[1])
