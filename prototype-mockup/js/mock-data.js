/**
 * FLAN MEDIA SERVER - PROTOTYPE MOCK DATA
 */

const FLAN_MOCK_DATA = {
  currentUser: null,

  // Whimsical Preset Avatars (High contrast, tactile SVG definitions)
  presetAvatars: {
    mascot: {
      name: "Anime Idol",
      svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100%" height="100%">
        <rect x="2" y="2" width="96" height="96" rx="6" fill="#3ea6ff" stroke="#000000" stroke-width="3"/>
        <rect x="6" y="6" width="88" height="88" rx="4" fill="#fdedd0"/>
        <polygon points="18,34 32,8 44,30" fill="#f37021" stroke="#000" stroke-width="2.5" stroke-linejoin="round"/>
        <polygon points="24,30 32,16 38,28" fill="#ffb4a2"/>
        <polygon points="56,30 68,8 82,34" fill="#f37021" stroke="#000" stroke-width="2.5" stroke-linejoin="round"/>
        <polygon points="62,28 68,16 76,30" fill="#ffb4a2"/>
        <path d="M16,48 C14,26 28,20 50,20 C72,20 86,26 84,48 C84,65 82,75 80,82 L70,82 C70,72 74,58 72,50 C70,44 65,48 60,52 C55,42 45,42 40,52 C35,48 30,44 28,50 C26,58 30,72 30,82 L20,82 C18,75 16,65 16,48 Z" fill="#f37021" stroke="#000" stroke-width="2.5" stroke-linejoin="round"/>
        <path d="M26,45 C26,70 34,78 50,78 C66,78 74,70 74,45 C74,32 66,28 50,28 C34,28 26,32 26,45 Z" fill="#fff5eb"/>
        <path d="M26,36 Q35,50 42,42 Q48,54 54,42 Q62,52 74,36 Q62,28 50,28 Q38,28 26,36 Z" fill="#f37021" stroke="#000" stroke-width="2" stroke-linejoin="round"/>
        <path d="M34,46 L42,51 L34,56" fill="none" stroke="#000" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M66,46 L58,51 L66,56" fill="none" stroke="#000" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
        <ellipse cx="32" cy="58" rx="4" ry="2.5" fill="#ff8a8a" opacity="0.8"/>
        <ellipse cx="68" cy="58" rx="4" ry="2.5" fill="#ff8a8a" opacity="0.8"/>
        <path d="M45,59 Q50,67 55,59 Z" fill="#d9383a" stroke="#000" stroke-width="1.8"/>
        <rect x="62" y="66" width="10" height="24" rx="3" transform="rotate(-25 67 78)" fill="#888" stroke="#000" stroke-width="2"/>
        <circle cx="63" cy="65" r="7" fill="#444" stroke="#000" stroke-width="2"/>
        <ellipse cx="58" cy="74" rx="6" ry="5" fill="#fff5eb" stroke="#000" stroke-width="2"/>
      </svg>`
    },
    flan: {
      name: "Flan Pudding",
      svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100%" height="100%">
        <rect x="2" y="2" width="96" height="96" rx="6" fill="#fbc02d" stroke="#000" stroke-width="3"/>
        <rect x="6" y="6" width="88" height="88" rx="4" fill="#fff9c4"/>
        <path d="M25,35 Q50,28 75,35 Q80,50 78,55 Q68,52 62,56 Q55,50 48,56 Q38,50 30,55 Q22,50 25,35 Z" fill="#6d4c41" stroke="#000" stroke-width="2"/>
        <polygon points="25,35 20,75 80,75 75,35" fill="#ffe082" stroke="#000" stroke-width="2.5" stroke-linejoin="round"/>
        <circle cx="40" cy="58" r="3.5" fill="#000"/>
        <circle cx="60" cy="58" r="3.5" fill="#000"/>
        <path d="M47,64 Q50,68 53,64" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round"/>
        <ellipse cx="34" cy="62" rx="3" ry="2" fill="#ff8a8a"/>
        <ellipse cx="66" cy="62" rx="3" ry="2" fill="#ff8a8a"/>
        <path d="M12,78 Q50,88 88,78" fill="none" stroke="#000" stroke-width="3" stroke-linecap="round"/>
      </svg>`
    },
    cat: {
      name: "Cool Cat",
      svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100%" height="100%">
        <rect x="2" y="2" width="96" height="96" rx="6" fill="#ab47bc" stroke="#000" stroke-width="3"/>
        <rect x="6" y="6" width="88" height="88" rx="4" fill="#f3e5f5"/>
        <polygon points="22,38 30,16 46,32" fill="#424242" stroke="#000" stroke-width="2.5"/>
        <polygon points="26,34 32,22 42,32" fill="#ff80ab"/>
        <polygon points="78,38 70,16 54,32" fill="#424242" stroke="#000" stroke-width="2.5"/>
        <polygon points="74,34 68,22 58,32" fill="#ff80ab"/>
        <ellipse cx="50" cy="52" rx="32" ry="26" fill="#424242" stroke="#000" stroke-width="2.5"/>
        <rect x="26" y="44" width="20" height="14" rx="2" fill="#ffe082" stroke="#000" stroke-width="2"/>
        <rect x="54" y="44" width="20" height="14" rx="2" fill="#ffe082" stroke="#000" stroke-width="2"/>
        <line x1="46" y1="51" x2="54" y2="51" stroke="#000" stroke-width="2.5"/>
        <polygon points="48,64 52,64 50,67" fill="#ff80ab"/>
        <path d="M47,68 Q50,71 53,68" fill="none" stroke="#fff" stroke-width="1.8"/>
        <line x1="16" y1="58" x2="30" y2="60" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
        <line x1="15" y1="66" x2="30" y2="64" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
        <line x1="84" y1="58" x2="70" y2="60" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
        <line x1="85" y1="66" x2="70" y2="64" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
      </svg>`
    },
    ghost: {
      name: "Friendly Ghost",
      svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100%" height="100%">
        <rect x="2" y="2" width="96" height="96" rx="6" fill="#26a69a" stroke="#000" stroke-width="3"/>
        <rect x="6" y="6" width="88" height="88" rx="4" fill="#e0f2f1"/>
        <path d="M25,50 C25,25 75,25 75,50 L75,76 Q68,70 61,76 Q54,70 50,76 Q46,70 39,76 Q32,70 25,76 Z" fill="#ffffff" stroke="#000" stroke-width="2.5" stroke-linejoin="round"/>
        <ellipse cx="40" cy="48" rx="4" ry="5.5" fill="#000"/>
        <ellipse cx="60" cy="48" rx="4" ry="5.5" fill="#000"/>
        <ellipse cx="34" cy="56" rx="3" ry="2" fill="#ff80ab"/>
        <ellipse cx="66" cy="56" rx="3" ry="2" fill="#ff80ab"/>
        <ellipse cx="50" cy="58" rx="3" ry="4" fill="#000"/>
      </svg>`
    },
    robot: {
      name: "Retro Bot",
      svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100%" height="100%">
        <rect x="2" y="2" width="96" height="96" rx="6" fill="#ff7043" stroke="#000" stroke-width="3"/>
        <rect x="6" y="6" width="88" height="88" rx="4" fill="#fbe9e7"/>
        <line x1="50" y1="26" x2="50" y2="16" stroke="#000" stroke-width="2.5"/>
        <circle cx="50" cy="14" r="4" fill="#ffeb3b" stroke="#000" stroke-width="2"/>
        <rect x="24" y="26" width="52" height="48" rx="6" fill="#b0bec5" stroke="#000" stroke-width="2.5"/>
        <rect x="32" y="36" width="36" height="18" rx="3" fill="#212121" stroke="#000" stroke-width="2"/>
        <circle cx="42" cy="45" r="3.5" fill="#00e676"/>
        <circle cx="58" cy="45" r="3.5" fill="#00e676"/>
        <rect x="38" y="60" width="24" height="6" rx="2" fill="#fff" stroke="#000" stroke-width="1.5"/>
        <line x1="44" y1="60" x2="44" y2="66" stroke="#000" stroke-width="1.5"/>
        <line x1="50" y1="60" x2="50" y2="66" stroke="#000" stroke-width="1.5"/>
        <line x1="56" y1="60" x2="56" y2="66" stroke="#000" stroke-width="1.5"/>
      </svg>`
    },
    star: {
      name: "Chubby Star",
      svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100%" height="100%">
        <rect x="2" y="2" width="96" height="96" rx="6" fill="#5c6bc0" stroke="#000" stroke-width="3"/>
        <rect x="6" y="6" width="88" height="88" rx="4" fill="#e8eaf6"/>
        <polygon points="50,14 60,37 85,38 65,54 72,78 50,64 28,78 35,54 15,38 40,37" fill="#ffd54f" stroke="#000" stroke-width="2.5" stroke-linejoin="round"/>
        <circle cx="44" cy="48" r="3" fill="#000"/>
        <circle cx="56" cy="48" r="3" fill="#000"/>
        <path d="M46,55 Q50,59 54,55" fill="none" stroke="#000" stroke-width="2" stroke-linecap="round"/>
        <ellipse cx="40" cy="53" rx="2.5" ry="1.5" fill="#ff8a8a"/>
        <ellipse cx="60" cy="53" rx="2.5" ry="1.5" fill="#ff8a8a"/>
      </svg>`
    }
  },

  users: [
    { id: 1, username: "mike", pin: "1234", role: "user", avatarType: "preset", avatarKey: "mascot", customAvatarData: null },
    { id: 2, username: "wesley", pin: "0000", role: "admin", avatarType: "preset", avatarKey: "flan", customAvatarData: null },
    { id: 3, username: "guest", pin: "1111", role: "user", avatarType: "preset", avatarKey: "ghost", customAvatarData: null }
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
        { id: 101, title: "S01E01 - Pilot", duration: "48m", progress: 100, positionSeconds: 2880, isFinished: true, streamUrl: "/stream/video/101", downloadUrl: "/download/video/101" },
        { id: 102, title: "S01E02 - Cat's in the Bag...", duration: "48m", progress: 50, positionSeconds: 1452, formattedPos: "24:12", isFinished: false, streamUrl: "/stream/video/102", downloadUrl: "/download/video/102" },
        { id: 103, title: "S01E03 - ...And the Bag's in the River", duration: "48m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/103", downloadUrl: "/download/video/103" },
        { id: 104, title: "S01E04 - Cancer Man", duration: "48m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/104", downloadUrl: "/download/video/104" }
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
        { id: 201, title: "Theatrical Cut (1080p)", duration: "2h 5m", progress: 75, positionSeconds: 5625, formattedPos: "1:33:45", isFinished: false, streamUrl: "/stream/video/201", downloadUrl: "/download/video/201" }
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
        { id: 301, title: "The Final Cut (2007)", duration: "1h 57m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/301", downloadUrl: "/download/video/301" },
        { id: 302, title: "Director's Cut (1992)", duration: "1h 56m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/302", downloadUrl: "/download/video/302" }
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
        { id: 401, title: "Session #1 - Asteroid Blues", duration: "24m", progress: 100, positionSeconds: 1440, isFinished: true, streamUrl: "/stream/video/401", downloadUrl: "/download/video/401" },
        { id: 402, title: "Session #2 - Stray Dog Strut", duration: "24m", progress: 20, positionSeconds: 288, formattedPos: "04:48", isFinished: false, streamUrl: "/stream/video/402", downloadUrl: "/download/video/402" }
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
        { id: 501, title: "Feature Film (Japanese Audio)", duration: "2h 14m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/501", downloadUrl: "/download/video/501" }
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
        { id: 1001, title: "Dune - Complete Edition", format: "epub", size: "2.4 MB", downloadUrl: "/download/book/1001" }
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
        { id: 1002, title: "Linux Programming Interface", format: "pdf", size: "28.5 MB", readUrl: "/stream/book/1002", downloadUrl: "/download/book/1002" }
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
        { id: 1003, title: "Neuromancer", format: "epub", size: "1.2 MB", downloadUrl: "/download/book/1003" }
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
