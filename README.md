# StayDB
StayDB is an in-memory storage engine written in [Go](https://golang.org/)
that uses compressed [bitmap indexing](https://en.wikipedia.org/wiki/Bitmap_index)
to implement fast set-like operations. Focused on supporting fast real time
analytics of large sets of data.
