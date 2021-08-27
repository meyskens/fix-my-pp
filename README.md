# Fix my PP

> Note: if you are not familiar with Beat Saber modding, do not attempt this!

A [recent change in BeatSaver broke PP display in song data on Beat Saber](https://github.com/halsafar/BeatSaberSongDataCore/issues/12). While a fix is on the way it is quite annoying as it breaks sorting of things. So I decided to write a fix of my own, while the scraper where the bug is in does not seem to be open source we can use the data file and the BeatSaver API to fix it ourselves. But no worries I do host a copy of my fix for you to use!

## Shut up, how do I fix it?
No worries I did most of the hard work for you already!

1) Download my version of [SongCoreData.dll](https://github.com/meyskens/fix-my-pp/raw/main/SongDataCore.dll) (don't trust me? then do it yourself!)
2) Go to your Beat Saber folder and go to `plugins` and replace SongCoreData.dll with mine
3) Repeat this every time you let ModAssistant replace it ;)

## I'm like you, a nerd, tell me more!
I wrote a small scripy in `main.go` that loads in the original JSON file and fetches the highest PP rank via scoresaber/beatsaver for that level then it places the PP score inside the data structure (as was suggested in the original issue!). Beatsaver almost got 90 000 songs so that takes a while. So I ran this on my laptop then hosted it just like the original. There is a tiny bug and that is that all levels got the standard gamemode PP as I didn't find the right integer translation yet (quite annoying that the API uses integers for all types but everything else does not).

To inject it into Beat Saber I modified BeatSaberSongDataCore to use my URL `https://static.eyskens.me/bssb/v2-all.json.gz`. Which hosts the modified file.
