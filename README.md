# Go-VMU
### Video Metadata Updater

## What Is It?

The Go Video Metadata Updater is designed to parse metadata from a Jellyfin created nfo file and update video files with that metadata.

The application:

Parses the nfo(xml)\
Extracts the metadata\
Updates the video files using ffmpeg and ffprobe

sample tags as returned by ffprobe after output:

```json
"tags": {
"title": "The Original",
"EPISODE": "1",
"ACTOR": "Evan Rachel Wood, Thandiwe Newton, Jeffrey Wright, Tessa Thompson, James Marsden, Angela Sarafyan, Anthony Hopkins, Ingrid Bolsø Berdal, Ed Harris, Shannon Woodward, Sidse Babett Knudsen, Jimmi Simpson, Ben Barnes, Simon Quarterman, Luke Hemsworth, Louis Herthum, Eddie Rouse, Kyle Bornheimer, Bridgid Coulter, Regi Davis, Mataeo Mingo, Trevante Rhodes, Micky Shiloah, Keller Wortham, Olivia May, Alex Marshall-Brown, Jeffrey Muller, Brook Kerr, Bradley Snedeker, Patrick Quinlan, David Coatsworth, Roberto Patino",
"TVDB_ID": "5626064",
"DIRECTOR": "Jonathan Nolan",
"WRITING_FRONTEND": "StaxRip v1.7.0.6",
"PLOT": "As another day of fantasy plays out in Westworld — a vast, remote park where guests pay top dollar to share wild-west adventures with android \"hosts\" — top programmer Bernard Lowe alerts park founder Dr. Robert Ford about incidents of aberrant behavior cropping up in some recently re-coded hosts. Meanwhile, in the Westworld town of Sweetwater, a rancher’s daughter named Dolores encounters a gunslinger named Teddy in the street — but their predictable narrative is upended by the appearance of a ruthless Man in Black and, later, by a supporting host’s unscripted encounter with an artifact of the outside world.",
"CREDITS": "Lisa Joy",
"GENRE": "Drama, Mystery, Sci-Fi",
"YEAR": "2016",
"WRITER": "Lisa Joy",
"SHOWTITLE": "Westworld",
"SEASON": "1",
"RUNTIME": "68",
"IMDB_ID": "tt4227538",
"ENCODER": "Lavf61.7.100"
}
```