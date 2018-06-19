This is a collection of various helper scripts I'd been keeping in my ~/bin folder and moving from computer to computer occasionally. I don't expect they will be useful to everybody, but one or two may help understand some usages of mplayer, ffmpeg, or the like. 

Or, if you're someone who likes code-reviewing some off-the-cuff bash, then this is for you!! Didn't reach for the best practices of functions or `local` and `readonly` too much, but I think I avoided the worst of bash, too.

2step
-----

If you're like me and don't have your phone constantly on you, or don't want a lot of 6-digit texts from Google, this script may be for you. It removes the phone from the "2-Step" authentication process, because when you use the Time-based One-Time Password (TOTP) method offered by things like the Google Authenticator app (whichever method that doesn't send you an SMS) it's not actually tied to your phone. You can easily generate TOTPs from your computer (or by hand if you're **really** quick at math). This script will do that. 

During the 2-step setup process, choose the non-SMS option and you will be presented with a code string that you're supposed to type/copy into Google Authenticator. Well, copy that code and put it in a file, `$HOME/.config/.2steps`, labeling it with whatever service you're logging into. E.g. put into the file: `gmail=sdfiu983rh0sd7fs07`. Now, to generate your 2step code to login, run the script, and give it the service name: `2step gmail` and it will print out the 6-digit code needed to login to your gmail, during the second step. See? No phone! This is why it is 2-step authentication and not 2-factor (SMS is closer, but I'm not satisfied there either).


bkup
----

Incomplete and non-functioning script I was going to use to backup my hard drives using JSON config. Manual config + cron won out in the end.


crypt
-----

Quick wrapper for using `gpg` to symmetrically encrypt (so, using AES instead of openPGP keys) files, and using a stdin pass prompt instead of the awful GUIs from `pinentry` (or curses-based).


fontview
--------

A shorthand script for passing the name of an installed font, and viewing it's charmap (XLFD not XFT). Use like: `fontview terminus`


gdate
----

Incomplete and non-functioning script I was using to correct git commit timestamps

ios-bkup
--------

basically just a long alias to check that proper HDs are mounted first, and run `idevicebackup2` with the right flags


lock
----

Will lock your X11 session using i3lock, and using a blurred, pixelated, annotated, or other desktop view based on parameters


lxc-setup
---------

Setup script I run before starting up an LXC container on an external drive. Checks to make sure drive is mounted with `dev` option, and sets up bridged networking.


mkgif
-----

Hopefully a more publicly useful script: converts a video file, using specified timestamps, into a gif automatically. It will even use youtube videos as video sources if you give it a URL instead of a filename. It's not /r/HighQualityGifs caliber output, but decent. Quality improvement suggestions are welcome (using imagemagick currently).

mkmov
-----

Another hopefully useful script to people: converts a bunch of still frames into a (selectable quality) movie. This time, the output can be rather stellar if you have the time and resources to wait. Given a folder full of files that may look like "frame_0001.jpg" through "frame_0860.jpg" in the `~/images` directory, the command will look like this: 

`mkmov ~/images/frame_%04d.jpg mymov.mkv`

This will create a pretty good timelapse video after a few minutes. The input file format gives the path to the images, the common part of all the file names (`frame_`) and where the numbers change, uses `%04d` to signal that the program should look for a 4-digit number there, preceded by 0's.  

Now let's say you messed up the start and end of the capture, and you only want to make the movie start from frame 120 and end at frame 760: `mkmov -s 120 -e 760 ~/images/frame_%04d.jpg mymov_shorter.mkv`

Now, our still frames were taken in 1080p, but we want our video to only be 720p, but be very good quality, and 60fps (meaning, the images will go by *very* quickly).

`mkmov -s 120 -e 760 -r 1280x720 -f 60 -q 9 ~/images/frame_%04d.jpg mymov_720.mkv`


Note that the higher the quality setting, the longer it will take to create the movie. This setting has a profound effect on the time taken to complete.

mnt
---

basically an alias script because I can never remember the flags and order to pass to `mount`, especially when assigning user rights flags on ntfs.


pixelate\_windows
----------------

Will pixelate your open X11 windows based on a screenshot of the current state. Produces an image where the background is untouched, windows are obscured. Based on py3lock.py by Airblader



notif-fix
---------

I think this had something to do with either notify-send or Chrome desktop notifications appearing in the wrong corner of my desktop, and I was trying to move them back. Not sure.

sc
---

This was an old script I used to "Set-Color" of my terminal, which didn't support an API for changing the theme/colors. It used the ANSI escape codes for changing the colors, and allowed putting presets in a colorconfig file.

sub
---

A small wrapper around sublime text to open normally when called the first time, but to spawn new windows when called anytime after that. The default sublime behavior when being called when already open is to just highlight the open window, instead of opening a new window.

tiling\_rules
------------

This is just the behavioral rules for my window manager; which windows should float vs. tile by default.

wacom
-----

userland driver for quick button remapping and orientation changing for wacom intuous pro drawing tablet

webcam
------

Small mplayer wrapper to view/stream a video source (e.g. liveview a webcam) and to optionally write it to stillframes for a slow-updating feed.

winfo
-----

Call this script then click on an X window to view some of it's properties (class, instance, title, type, state)


License
=======

All code in this repository is licensed under the MIT license (see `LICENSE` file) and  
Copyright (c) 2015 Dan Panzarella
