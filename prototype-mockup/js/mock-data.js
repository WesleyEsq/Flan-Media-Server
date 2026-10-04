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
        <rect x="2" y="2" width="96" height="96" rx="6" fill="#fdedd0"/>
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
      title: "Breaking Bad",
      releaseYear: 2008,
      type: "series",
      sourceId: 1,
      lastWatched: "2026-10-03T21:15:00Z",
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
      title: "Spirited Away",
      releaseYear: 2001,
      type: "movie",
      sourceId: 2,
      lastWatched: "2026-10-02T19:40:00Z",
      overview: "During her family's move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches and spirits, and where humans are changed into beasts.",
      coverColor: "#e65100",
      badge: "Movie",
      files: [
        { id: 201, title: "Theatrical Cut (1080p)", duration: "2h 5m", progress: 75, positionSeconds: 5625, formattedPos: "1:33:45", isFinished: false, streamUrl: "/stream/video/201", downloadUrl: "/download/video/201" }
      ]
    },
    {
      id: 3,
      title: "Blade Runner",
      releaseYear: 1982,
      type: "movie",
      sourceId: 1,
      lastWatched: null,
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
      title: "Cowboy Bebop",
      releaseYear: 1998,
      type: "series",
      sourceId: 1,
      lastWatched: "2026-10-01T22:10:00Z",
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
      title: "Princess Mononoke",
      releaseYear: 1997,
      type: "movie",
      sourceId: 2,
      lastWatched: null,
      overview: "On a journey to find the cure for a Tatarigami's curse, Ashitaka finds himself in the middle of a war between the forest gods and Tatara, a mining colony.",
      coverColor: "#00695c",
      badge: "Movie",
      files: [
        { id: 501, title: "Feature Film (Japanese Audio)", duration: "2h 14m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/501", downloadUrl: "/download/video/501" }
      ]
    },
    {
      id: 6,
      title: "Alien",
      releaseYear: 1979,
      type: "movie",
      sourceId: 1,
      lastWatched: null,
      overview: "The crew of a commercial spacecraft encounters a deadly lifeform after investigating an unknown transmission.",
      coverColor: "#263238",
      badge: "Movie",
      files: [
        { id: 601, title: "Director's Cut (1080p)", duration: "1h 56m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/601", downloadUrl: "/download/video/601" }
      ]
    },
    {
      id: 7,
      title: "Akira",
      releaseYear: 1988,
      type: "movie",
      sourceId: 2,
      lastWatched: null,
      overview: "A secret military project endangers Neo-Tokyo when it turns a biker gang member into a rampaging psychic psychopath.",
      coverColor: "#b71c1c",
      badge: "Movie",
      files: [
        { id: 701, title: "Remastered (1080p)", duration: "2h 4m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/701", downloadUrl: "/download/video/701" }
      ]
    },
    {
      id: 8,
      title: "The Wire",
      releaseYear: 2002,
      type: "series",
      sourceId: 2,
      lastWatched: null,
      overview: "Told from the points of view of both the Baltimore police and their targets, the series captures a universe where easy distinctions between good and evil are discarded.",
      coverColor: "#37474f",
      badge: "Series",
      files: [
        { id: 801, title: "S01E01 - The Target", duration: "1h 0m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/801", downloadUrl: "/download/video/801" },
        { id: 802, title: "S01E02 - The Detail", duration: "56m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/802", downloadUrl: "/download/video/802" }
      ]
    },
    {
      id: 9,
      title: "Dune",
      releaseYear: 2021,
      type: "movie",
      sourceId: 2,
      lastWatched: null,
      overview: "A noble family becomes embroiled in a war for control over the galaxy's most valuable asset while its heir becomes troubled by visions of a dark future.",
      coverColor: "#d84315",
      badge: "Movie",
      files: [
        { id: 901, title: "Part One (4K HDR)", duration: "2h 35m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/901", downloadUrl: "/download/video/901" }
      ]
    },
    {
      id: 10,
      title: "Severance",
      releaseYear: 2022,
      type: "series",
      sourceId: 1,
      lastWatched: null,
      overview: "Mark leads a team of office workers whose memories have been surgically divided between their work and personal lives.",
      coverColor: "#0277bd",
      badge: "Series",
      files: [
        { id: 1001, title: "S01E01 - Good News About Hell", duration: "57m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1001", downloadUrl: "/download/video/1001" }
      ]
    },
    {
      id: 11,
      title: "The Matrix",
      releaseYear: 1999,
      type: "movie",
      sourceId: 1,
      lastWatched: null,
      overview: "A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.",
      coverColor: "#1b5e20",
      badge: "Movie",
      files: [
        { id: 1101, title: "Theatrical Cut (1080p)", duration: "2h 16m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1101", downloadUrl: "/download/video/1101" }
      ]
    },
    {
      id: 12,
      title: "Cyberpunk: Edgerunners",
      releaseYear: 2022,
      type: "series",
      sourceId: 2,
      lastWatched: null,
      overview: "A street kid trying to survive in a technology and body modification-obsessed city of the future.",
      coverColor: "#c2185b",
      badge: "Series",
      files: [
        { id: 1201, title: "Episode 1 - Let You Down", duration: "24m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1201", downloadUrl: "/download/video/1201" }
      ]
    },
    {
      id: 13,
      title: "Fargo",
      releaseYear: 1996,
      type: "movie",
      sourceId: 1,
      lastWatched: null,
      overview: "Minnesota police chief Marge Gunderson investigates homicides that occurred after a car salesman hired two criminals to kidnap his wife.",
      coverColor: "#455a64",
      badge: "Movie",
      files: [
        { id: 1301, title: "Feature Film", duration: "1h 38m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1301", downloadUrl: "/download/video/1301" }
      ]
    },
    {
      id: 14,
      title: "Metropolis",
      releaseYear: 1927,
      type: "movie",
      sourceId: 1,
      lastWatched: null,
      overview: "In a futuristic city sharply divided between the working class and the city planners, the son of the city's mastermind falls in love with a working-class prophet.",
      coverColor: "#3e2723",
      badge: "Movie",
      files: [
        { id: 1401, title: "Restored Giorgio Moroder Edition", duration: "2h 33m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1401", downloadUrl: "/download/video/1401" }
      ]
    },
    {
      id: 15,
      title: "Succession",
      releaseYear: 2018,
      type: "series",
      sourceId: 2,
      lastWatched: null,
      overview: "The Roy family is known for controlling the biggest media and entertainment company in the world. However, their world changes when their aging father steps down.",
      coverColor: "#424242",
      badge: "Series",
      files: [
        { id: 1501, title: "S01E01 - Celebration", duration: "1h 1m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1501", downloadUrl: "/download/video/1501" }
      ]
    },
    {
      id: 16,
      title: "Twin Peaks",
      releaseYear: 1990,
      type: "series",
      sourceId: 1,
      lastWatched: null,
      overview: "An idiosyncratic FBI agent investigates the murder of a young woman in the even more idiosyncratic town of Twin Peaks.",
      coverColor: "#4e342e",
      badge: "Series",
      files: [
        { id: 1601, title: "Pilot - Northwest Passage", duration: "1h 34m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1601", downloadUrl: "/download/video/1601" }
      ]
    },
    {
      id: 17,
      title: "Neon Genesis Evangelion",
      releaseYear: 1995,
      type: "series",
      sourceId: 1,
      lastWatched: null,
      overview: "A teenage boy finds himself recruited by his estranged father into the shadowy organization NERV to pilot a giant bio-machine.",
      coverColor: "#6a1b9a",
      badge: "Series",
      files: [
        { id: 1701, title: "Episode 1 - Angel Attack", duration: "24m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1701", downloadUrl: "/download/video/1701" }
      ]
    },
    {
      id: 18,
      title: "Chernobyl",
      releaseYear: 2019,
      type: "series",
      sourceId: 2,
      lastWatched: null,
      overview: "In April 1986, an explosion at the Chernobyl nuclear power plant becomes one of the world's worst man-made catastrophes.",
      coverColor: "#f57f17",
      badge: "Series",
      files: [
        { id: 1801, title: "1:23:45", duration: "59m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1801", downloadUrl: "/download/video/1801" }
      ]
    },
    {
      id: 19,
      title: "Mad Max: Fury Road",
      releaseYear: 2015,
      type: "movie",
      sourceId: 2,
      lastWatched: null,
      overview: "In a post-apocalyptic wasteland, a woman rebels against a tyrannical ruler in search for her homeland with the aid of a group of female prisoners.",
      coverColor: "#e65100",
      badge: "Movie",
      files: [
        { id: 1901, title: "Black & Chrome Edition", duration: "2h 0m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/1901", downloadUrl: "/download/video/1901" }
      ]
    },
    {
      id: 20,
      title: "Better Call Saul",
      releaseYear: 2015,
      type: "series",
      sourceId: 1,
      lastWatched: null,
      overview: "The trials and tribulations of criminal lawyer Jimmy McGill in the years leading up to his fateful run-in with Walter White and Jesse Pinkman.",
      coverColor: "#f9a825",
      badge: "Series",
      files: [
        { id: 2001, title: "S01E01 - Uno", duration: "53m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/2001", downloadUrl: "/download/video/2001" }
      ]
    },
    {
      id: 21,
      title: "Ghost in the Shell",
      releaseYear: 1995,
      type: "movie",
      sourceId: 2,
      lastWatched: null,
      overview: "A cyborg policewoman and her partner hunt a mysterious and powerful hacker called the Puppet Master.",
      coverColor: "#00838f",
      badge: "Movie",
      files: [
        { id: 2101, title: "Original 1995 Theatrical (1080p)", duration: "1h 23m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/2101", downloadUrl: "/download/video/2101" }
      ]
    },
    {
      id: 22,
      title: "Taxi Driver",
      releaseYear: 1976,
      type: "movie",
      sourceId: 1,
      lastWatched: null,
      overview: "A mentally unstable veteran works as a nighttime taxi driver in New York City, where the perceived decadence fuels his urge for violent action.",
      coverColor: "#d32f2f",
      badge: "Movie",
      files: [
        { id: 2201, title: "4K Master (1080p)", duration: "1h 54m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/2201", downloadUrl: "/download/video/2201" }
      ]
    },
    {
      id: 23,
      title: "Blade Runner 2049",
      releaseYear: 2017,
      type: "movie",
      sourceId: 2,
      lastWatched: null,
      overview: "Young Blade Runner K's discovery of a long-buried secret leads him to track down former Blade Runner Rick Deckard, who's been missing for thirty years.",
      coverColor: "#0288d1",
      badge: "Movie",
      files: [
        { id: 2301, title: "Theatrical Release (1080p)", duration: "2h 44m", progress: 0, positionSeconds: 0, isFinished: false, streamUrl: "/stream/video/2301", downloadUrl: "/download/video/2301" }
      ]
    }
  ],

  books: [
    {
      id: 101,
      title: "Dune",
      author: "Frank Herbert",
      releaseYear: 1965,
      sourceId: 3,
      status: "reading",
      lastRead: "2026-10-03T18:00:00Z",
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
      releaseYear: 2010,
      sourceId: 3,
      status: "unread",
      lastRead: null,
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
      releaseYear: 1984,
      sourceId: 3,
      status: "finished",
      lastRead: null,
      overview: "Case was the sharpest data-thief in the business, until he pissed off the wrong people and they crippled his nervous system.",
      coverColor: "#4a148c",
      badge: "EPUB",
      files: [
        { id: 1003, title: "Neuromancer", format: "epub", size: "1.2 MB", downloadUrl: "/download/book/1003" }
      ]
    },
    {
      id: 104,
      title: "Snow Crash",
      author: "Neal Stephenson",
      releaseYear: 1992,
      sourceId: 3,
      status: "reading",
      lastRead: "2026-10-02T14:30:00Z",
      overview: "In reality, Hiro Protagonist delivers pizza for Uncle Enzo's CosoNostra Pizza Inc., but in the Metaverse he's a warrior prince.",
      coverColor: "#00695c",
      badge: "EPUB",
      files: [
        { id: 1004, title: "Snow Crash - Complete", format: "epub", size: "1.8 MB", downloadUrl: "/download/book/1004" }
      ]
    },
    {
      id: 105,
      title: "Designing Data-Intensive Applications",
      author: "Martin Kleppmann",
      releaseYear: 2017,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "Data is at the center of many challenges in system design today. Difficult issues need to be figured out, such as scalability, consistency, reliability, efficiency, and maintainability.",
      coverColor: "#004d40",
      badge: "PDF",
      files: [
        { id: 1005, title: "Designing Data-Intensive Applications (1st Ed)", format: "pdf", size: "14.2 MB", readUrl: "/stream/book/1005", downloadUrl: "/download/book/1005" }
      ]
    },
    {
      id: 106,
      title: "The Pragmatic Programmer",
      author: "Andrew Hunt, David Thomas",
      releaseYear: 1999,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "Illustrates the best approaches and major pitfalls of many different aspects of software development.",
      coverColor: "#bf360c",
      badge: "EPUB",
      files: [
        { id: 1006, title: "The Pragmatic Programmer (20th Anniv)", format: "epub", size: "3.1 MB", downloadUrl: "/download/book/1006" }
      ]
    },
    {
      id: 107,
      title: "Clean Architecture",
      author: "Robert C. Martin",
      releaseYear: 2017,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "By applying universal rules of software architecture, you can dramatically improve developer productivity throughout the life of any software system.",
      coverColor: "#1a237e",
      badge: "PDF",
      files: [
        { id: 1007, title: "Clean Architecture", format: "pdf", size: "9.8 MB", readUrl: "/stream/book/1007", downloadUrl: "/download/book/1007" }
      ]
    },
    {
      id: 108,
      title: "Structure and Interpretation of Computer Programs",
      author: "Harold Abelson, Gerald Jay Sussman",
      releaseYear: 1984,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "SICP has had a dramatic impact on computer science curricula over the past decades, exploring functional programming and metacircular evaluators.",
      coverColor: "#311b92",
      badge: "PDF",
      files: [
        { id: 1008, title: "SICP (2nd Edition)", format: "pdf", size: "18.4 MB", readUrl: "/stream/book/1008", downloadUrl: "/download/book/1008" }
      ]
    },
    {
      id: 109,
      title: "Godel, Escher, Bach",
      author: "Douglas Hofstadter",
      releaseYear: 1979,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "An Eternal Golden Braid: A metaphorical fugue on minds and machines in the spirit of Lewis Carroll.",
      coverColor: "#4e342e",
      badge: "PDF",
      files: [
        { id: 1009, title: "Godel, Escher, Bach - Complete Scan", format: "pdf", size: "32.1 MB", readUrl: "/stream/book/1009", downloadUrl: "/download/book/1009" }
      ]
    },
    {
      id: 110,
      title: "The Mythical Man-Month",
      author: "Fred Brooks",
      releaseYear: 1975,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "Essays on software engineering and project management: Adding manpower to a late software project makes it later.",
      coverColor: "#263238",
      badge: "EPUB",
      files: [
        { id: 1010, title: "The Mythical Man-Month (Anniversary Ed)", format: "epub", size: "2.1 MB", downloadUrl: "/download/book/1010" }
      ]
    },
    {
      id: 111,
      title: "Foundation",
      author: "Isaac Asimov",
      releaseYear: 1951,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "The story of our civilization, the Galactic Empire, now in its twilight, and Hari Seldon's mathematical science of psychohistory.",
      coverColor: "#006064",
      badge: "EPUB",
      files: [
        { id: 1011, title: "Foundation - Volume 1", format: "epub", size: "1.5 MB", downloadUrl: "/download/book/1011" }
      ]
    },
    {
      id: 112,
      title: "Solaris",
      author: "Stanislaw Lem",
      releaseYear: 1961,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "A psychologist arrives at a research station hovering above the oceanic surface of the planet Solaris, studying its sentient alien ocean.",
      coverColor: "#01579b",
      badge: "EPUB",
      files: [
        { id: 1012, title: "Solaris (Bilingual Translation)", format: "epub", size: "1.9 MB", downloadUrl: "/download/book/1012" }
      ]
    },
    {
      id: 113,
      title: "Do Androids Dream of Electric Sheep?",
      author: "Philip K. Dick",
      releaseYear: 1968,
      sourceId: 3,
      status: "unread",
      lastRead: null,
      overview: "By 2021, the World War has killed millions, driving entire species into extinction and sending mankind off-planet. Rick Deckard is tasked with retiring rogue Nexus-6 androids.",
      coverColor: "#3e2723",
      badge: "EPUB",
      files: [
        { id: 1013, title: "Do Androids Dream of Electric Sheep?", format: "epub", size: "1.7 MB", downloadUrl: "/download/book/1013" }
      ]
    }
  ],

  serverMetrics: {
    heapUsed: "18.5 MB",
    sysMem: "28.2 MB",
    goroutines: 14,
    activeStreams: 1,
    maxStreams: 3,
    storageMediaFree: "238 GB / 500 GB (External HDD)",
    storageDbFree: "24.8 GB / 32 GB (eMMC / SD)",
    uptime: "14 days, 3 hours"
  },

  storageSources: [
    { id: 1, name: "Primary Video Drive", mediaType: "video", path: "./media/video", freeSpace: "184 GB / 500 GB", isActive: true },
    { id: 2, name: "4TB USB Movie Drive", mediaType: "video", path: "/mnt/usb1/movies", freeSpace: "2.1 TB / 4.0 TB", isActive: true },
    { id: 3, name: "Household Books", mediaType: "book", path: "./media/books", freeSpace: "58 GB / 128 GB", isActive: true }
  ],

  pendingIngestion: [
    {
      id: 901,
      sourceId: 2,
      sourceName: "4TB USB Movie Drive",
      rawFolder: "the.wire.s01.720p.hdtv.x264-ctu[eztv]",
      proposedTitle: "The Wire",
      proposedYear: 2002,
      mediaType: "video",
      fileCount: 13,
      coverColor: "#b71c1c",
      badge: "Series",
      overview: "Told from the points of view of both the Baltimore police and their targets, the series captures a universe where easy distinctions between good and evil are discarded.",
      files: [
        { id: 9001, rawFilename: "the.wire.s01e01.720p.mkv", title: "S01E01 - The Target", duration: "60m", progress: 0, positionSeconds: 0, isFinished: false, isHidden: false, streamUrl: "#", downloadUrl: "#" },
        { id: 9002, rawFilename: "the.wire.s01e02.720p.mkv", title: "S01E02 - The Detail", duration: "56m", progress: 0, positionSeconds: 0, isFinished: false, isHidden: false, streamUrl: "#", downloadUrl: "#" },
        { id: 9003, rawFilename: "the.wire.s01e03.720p.mkv", title: "S01E03 - The Buys", duration: "55m", progress: 0, positionSeconds: 0, isFinished: false, isHidden: false, streamUrl: "#", downloadUrl: "#" }
      ]
    },
    {
      id: 902,
      sourceId: 1,
      sourceName: "Primary Video Drive",
      rawFolder: "alien.1979.directors.cut.1080p.bluray.x264",
      proposedTitle: "Alien",
      proposedYear: 1979,
      mediaType: "video",
      fileCount: 2,
      coverColor: "#263238",
      badge: "Movie",
      overview: "The crew of a commercial spacecraft encounters a deadly lifeform after investigating an unknown transmission.",
      files: [
        { id: 9010, rawFilename: "Alien.1979.Directors.Cut.1080p.mkv", title: "Director's Cut", duration: "1h 56m", progress: 0, positionSeconds: 0, isFinished: false, isHidden: false, streamUrl: "#", downloadUrl: "#" },
        { id: 9011, rawFilename: "Alien.1979.Theatrical.1080p.mkv", title: "Theatrical Cut", duration: "1h 57m", progress: 0, positionSeconds: 0, isFinished: false, isHidden: false, streamUrl: "#", downloadUrl: "#" }
      ]
    },
    {
      id: 903,
      sourceId: 3,
      sourceName: "Household Books",
      rawFolder: "Snow.Crash.Neal.Stephenson.1992.epub",
      proposedTitle: "Snow Crash",
      proposedAuthor: "Neal Stephenson",
      proposedYear: 1992,
      mediaType: "book",
      fileCount: 1,
      coverColor: "#00695c",
      badge: "EPUB",
      overview: "In reality, Hiro Protagonist delivers pizza for Uncle Enzo's CosoNostra Pizza Inc., but in the Metaverse he's a warrior prince.",
      files: [
        { id: 9020, rawFilename: "Snow Crash - Neal Stephenson.epub", title: "Snow Crash", format: "epub", size: "1.8 MB", downloadUrl: "#" }
      ]
    }
  ]
};

