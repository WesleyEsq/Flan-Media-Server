/**
 * FLAN MEDIA SERVER - PROTOTYPE MOCK DATA
 */

const FLAN_MOCK_DATA = {
  currentUser: null,

  users: [
    { id: 1, username: "mike", pin: "1234", role: "user", avatar: "mascot" },
    { id: 2, username: "wesley", pin: "0000", role: "admin", avatar: "mascot" },
    { id: 3, username: "guest", pin: "1111", role: "user", avatar: "default" }
  ],

  videos: [
    {
      id: 1,
      title: "Breaking Bad (2008)",
      type: "series",
      overview: "A high school chemistry teacher diagnosed with inoperable lung cancer turns to manufacturing and selling methamphetamine with a former student.",
      coverColor: "#2e7d32",
      badge: "Series",
      files: [
        { id: 101, title: "S01E01 - Pilot", duration: "48m", progress: 100, positionSeconds: 2880, isFinished: true },
        { id: 102, title: "S01E02 - Cat's in the Bag...", duration: "48m", progress: 50, positionSeconds: 1452, formattedPos: "24:12", isFinished: false },
        { id: 103, title: "S01E03 - ...And the Bag's in the River", duration: "48m", progress: 0, positionSeconds: 0, isFinished: false },
        { id: 104, title: "S01E04 - Cancer Man", duration: "48m", progress: 0, positionSeconds: 0, isFinished: false }
      ]
    },
    {
      id: 2,
      title: "Spirited Away (2001)",
      type: "movie",
      overview: "During her family's move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches and spirits, and where humans are changed into beasts.",
      coverColor: "#e65100",
      badge: "Movie",
      files: [
        { id: 201, title: "Theatrical Cut (1080p)", duration: "2h 5m", progress: 75, positionSeconds: 5625, formattedPos: "1:33:45", isFinished: false }
      ]
    },
    {
      id: 3,
      title: "Blade Runner (1982)",
      type: "movie",
      overview: "A blade runner must pursue and terminate four replicants who stole a ship in space and have returned to Earth to find their creator.",
      coverColor: "#1565c0",
      badge: "Movie",
      files: [
        { id: 301, title: "The Final Cut (2007)", duration: "1h 57m", progress: 0, positionSeconds: 0, isFinished: false },
        { id: 302, title: "Director's Cut (1992)", duration: "1h 56m", progress: 0, positionSeconds: 0, isFinished: false }
      ]
    },
    {
      id: 4,
      title: "Cowboy Bebop (1998)",
      type: "series",
      overview: "The futuristic misadventures and tragedies of an easygoing bounty hunter and his partners.",
      coverColor: "#c2185b",
      badge: "Series",
      files: [
        { id: 401, title: "Session #1 - Asteroid Blues", duration: "24m", progress: 100, positionSeconds: 1440, isFinished: true },
        { id: 402, title: "Session #2 - Stray Dog Strut", duration: "24m", progress: 20, positionSeconds: 288, formattedPos: "04:48", isFinished: false }
      ]
    },
    {
      id: 5,
      title: "Princess Mononoke (1997)",
      type: "movie",
      overview: "On a journey to find the cure for a Tatarigami's curse, Ashitaka finds himself in the middle of a war between the forest gods and Tatara, a mining colony.",
      coverColor: "#00695c",
      badge: "Movie",
      files: [
        { id: 501, title: "Feature Film (Japanese Audio)", duration: "2h 14m", progress: 0, positionSeconds: 0, isFinished: false }
      ]
    }
  ],

  books: [
    {
      id: 101,
      title: "Dune",
      author: "Frank Herbert",
      overview: "Set on the desert planet Arrakis, Dune is the story of Paul Atreides—who would become known as Muad'Dib—and of a great family's ambition to bring to fruition humankind's most ancient and unattainable dream.",
      coverColor: "#d84315",
      badge: "EPUB",
      files: [
        { id: 1001, title: "Dune - Complete Edition (EPUB)", size: "2.4 MB", progress: 34, page: "Chapter 12" }
      ]
    },
    {
      id: 102,
      title: "The Linux Programming Interface",
      author: "Michael Kerrisk",
      overview: "The definitive guide to the Linux and UNIX programming interface—the interface employed by nearly every application that runs on a Linux or UNIX system.",
      coverColor: "#283593",
      badge: "PDF",
      files: [
        { id: 1002, title: "Linux Programming Interface.pdf", size: "28.5 MB", progress: 15, page: "Page 234 / 1552" }
      ]
    },
    {
      id: 103,
      title: "Neuromancer",
      author: "William Gibson",
      overview: "Case was the sharpest data-thief in the business, until he pissed off the wrong people and they crippled his nervous system.",
      coverColor: "#4a148c",
      badge: "EPUB",
      files: [
        { id: 1003, title: "Neuromancer.epub", size: "1.2 MB", progress: 0, page: "Start" }
      ]
    }
  ],

  serverMetrics: {
    memoryUsed: "14.2 MB",
    memoryLimit: "16.0 MB",
    activeStreams: 1,
    maxStreams: 3,
    storageMediaFree: "238 GB / 500 GB (External HDD)",
    storageDbFree: "24.8 GB / 32 GB (eMMC / SD)",
    uptime: "14 days, 3 hours"
  }
};
