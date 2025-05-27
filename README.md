# Go-VMU
### Video Metadata Updater

## What Is It?

The Go Video Metadata Updater (Go-VMU) is designed to parse metadata from Jellyfin-created NFO files and update video files with that metadata. This tool helps maintain consistent metadata across your media library by extracting information from NFO files and embedding it directly into the video files.

## How It Works

The application follows a simple workflow:

1. Parses the NFO (XML) file associated with a video file
2. Extracts the metadata (title, plot, actors, etc.)
3. Updates the video files using FFmpeg while preserving the original video and audio quality
4. Validates the updated file to ensure no corruption occurred

## Requirements

- Go 1.16 or higher
- FFmpeg and FFprobe installed and available in your PATH
- Jellyfin-compatible NFO files (Recommended to use Jellyfin automated nfo creation)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-vmu.git
cd go-vmu

# Build the application
go build -o vmu ./cmd/vmu
```

## Configuration

Create a `config.toml` file in the application directory or use the provided `sample-config.toml` as a template:

```toml
nfo_type = "jellyfin"
scan_folder = "/path/to/your/media/library"

[logger]
level = "info"        # debug, info, warn, error
pretty = true         # Pretty print logs
time_format = 1979-05-27T07:32:00
log_file = "./vmu.log"
max_size = 4          # Max size in MB
max_age = 14          # Max age in days
max_backups = 10      # Max number of backup files
compress_backup = false
```

## Usage

Currently, the application is in development and the main.go file contains an example implementation. To use the application:

1. Update the configuration file with your media library path
2. Modify the main.go file to process your desired files or implement a file walker
3. Run the application:

```bash
go run cmd/vmu/main.go
```

### Future Enhancements

- Command-line interface for easier usage
- Recursive directory scanning
- Concurrent processing for faster updates
- Support for different NFO formats
- Batch processing capabilities

## Sample Output

After processing, FFprobe will show the embedded metadata tags in your video files:

### Before
```json
{
  "tags": {
    "encoder": "libebml v1.3.5 + libmatroska v1.4.8",
    "creation_time": "2020-03-04T08:14:33.000000Z",
    "Writing frontend": "StaxRip v1.7.0.6"
  }
}
```

### After
```json
{
  "tags": {
    "title": "The Original",
    "EPISODE": "1",
    "ACTOR": "Evan Rachel Wood, Thandiwe Newton, Jeffrey Wright, Tessa Thompson, James Marsden, Angela Sarafyan, Anthony Hopkins, Ingrid Bolsø Berdal, Ed Harris, Shannon Woodward, Sidse Babett Knudsen, Jimmi Simpson, Ben Barnes, Simon Quarterman, Luke Hemsworth, Louis Herthum, Eddie Rouse, Kyle Bornheimer, Bridgid Coulter, Regi Davis, Mataeo Mingo, Trevante Rhodes, Micky Shiloah, Keller Wortham, Olivia May, Alex Marshall-Brown, Jeffrey Muller, Brook Kerr, Bradley Snedeker, Patrick Quinlan, David Coatsworth, Roberto Patino",
    "TVDB_ID": "5626064",
    "DIRECTOR": "Jonathan Nolan",
    "WRITING_FRONTEND": "StaxRip v1.7.0.6",
    "PLOT": "As another day of fantasy plays out in Westworld — a vast, remote park where guests pay top dollar to share wild-west adventures with android \"hosts\" — top programmer Bernard Lowe alerts park founder Dr. Robert Ford about incidents of aberrant behavior cropping up in some recently re-coded hosts. Meanwhile, in the Westworld town of Sweetwater, a rancher's daughter named Dolores encounters a gunslinger named Teddy in the street — but their predictable narrative is upended by the appearance of a ruthless Man in Black and, later, by a supporting host's unscripted encounter with an artifact of the outside world.",
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
}
```

## Limitations and Known Issues

- NFS shares may cause issues with file operations
- Processing large libraries can be time-consuming without concurrency
- Currently only supports Jellyfin NFO format
- The application creates temporary files during processing which require additional disk space
