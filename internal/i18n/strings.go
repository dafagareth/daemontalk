package i18n

type UI struct {
	Nav_Home     string
	Nav_Behind   string
	Nav_Terminal string
	Footer_Copy  string

	Behind_Title        string
	Behind_Colophon     string
	Behind_ColophonBody string
	Behind_Uses         string
	Behind_Learning     string

	Footer_Bio       string
	Footer_Read      string
	Footer_Site      string
	Footer_Topics    string
	Footer_Lab       string
	Footer_Resources string
	Footer_Meta      string
	Footer_Legal     string
	Footer_Community string
	Footer_Ecosystem string
	Footer_Other     string
	Privacy_Title    string

	Home_BehindDoor string

	Blog_Title             string
	Blog_ReadMore          string
	Blog_MinRead           string
	Blog_NoPosts           string
	Blog_SearchPlaceholder string
	Blog_NoResults         string
	Blog_PrevPage          string
	Blog_NextPage          string
	Blog_LoadMore          string

	Post_Tags    string
	Post_Back    string
	Post_Related string
	Post_TOC     string

	Comments_Heading    string
	Comments_Empty      string
	Comments_NameLabel  string
	Comments_BodyLabel  string
	Comments_Submit     string
	Comments_Error      string
	Comments_Delete     string
	Comments_Reply      string
	Comments_Cancel     string
	Comments_ReplyingTo string

	Post_Views string

	About_Intro    string
	About_Building string
	About_Stack    string
	About_Contact  string

	Contact_FormTitle    string
	Contact_NameLabel    string
	Contact_EmailLabel   string
	Contact_MessageLabel string
	Contact_Submit       string
	Contact_Success      string
	Contact_Error        string

	NotFound_Title string
	NotFound_Body  string
	NotFound_Home  string

	Forbidden_Title string
	Forbidden_Body  string
	Forbidden_Home  string

	ServerError_Title string
	ServerError_Body  string
	ServerError_Home  string

	Error_SearchPlaceholder string
	Error_Diagnostics       string

	Uses_Hardware string
	Uses_Software string
	Uses_Desk     string

	Now_WorkingOn string
	Now_Reading   string
	Now_Learning  string

	Reading_Title    string
	Reading_Subtitle string

	Stats_PageTitle string

	Saved_Title    string
	Saved_Subtitle string
	Saved_Empty    string
	Saved_Clear    string

	Lang_Switch string

	Search_Title      string
	Search_ResultsFor string

	Guestbook_Title     string
	Guestbook_NameLabel string
	Guestbook_BodyLabel string
	Guestbook_Submit    string
	Guestbook_Empty     string

	Stats_Title string

	Changelog_Title string
	Links_Title     string
	Links_Subtitle  string

	Terms_Title       string
	Terms_LastUpdated string
}

var Strings = map[string]UI{
	"en": {
		Nav_Home:     "Home",
		Nav_Behind:   "behind",
		Nav_Terminal: "terminal",
		Footer_Copy:  "© 2026 Dafa",

		Behind_Title:        "Behind this website",
		Behind_Colophon:     "Colophon",
		Behind_ColophonBody: "This site is a single Go binary: chi for routing, templ for HTML templates, HTMX for the few interactive bits, Tailwind for styling, SQLite for comments and views. Server-side rendered, self-hosted, zero trackers.",
		Behind_Uses:         "Stack & uses",
		Behind_Learning:     "Now / learning",

		Footer_Bio:       "Open engineering notebook and learning log. Exploring Go backend, Python, Linux systems, and distributed architecture.",
		Footer_Read:      "Navigation",
		Footer_Site:      "Site",
		Footer_Topics:    "Topics",
		Footer_Lab:       "Lab",
		Footer_Resources: "Resources",
		Footer_Meta:      "Meta",
		Footer_Legal:     "Legal & Policy",
		Footer_Community: "Open Source",
		Footer_Ecosystem: "Feeds & Meta",
		Footer_Other:     "Other & Tools",
		Privacy_Title:    "Privacy Policy",

		Home_BehindDoor: "look behind this website",

		Blog_Title:             "Blog",
		Blog_ReadMore:          "Read",
		Blog_MinRead:           "min read",
		Blog_NoPosts:           "No posts yet.",
		Blog_SearchPlaceholder: "Search posts...",
		Blog_NoResults:         "No posts match your search.",
		Blog_PrevPage:          "Newer",
		Blog_NextPage:          "Older",
		Blog_LoadMore:          "Load more",

		Post_Tags:    "Tags",
		Post_Back:    "Back to blog",
		Post_Related: "Related posts",
		Post_TOC:     "Table of Contents",

		Comments_Heading:    "Comments",
		Comments_Empty:      "No comments yet. Be the first to comment!",
		Comments_NameLabel:  "Your name",
		Comments_BodyLabel:  "Write a comment...",
		Comments_Submit:     "Submit comment",
		Comments_Error:      "Failed to post comment. Make sure name and body are filled.",
		Comments_Delete:     "Delete",
		Comments_Reply:      "Reply",
		Comments_Cancel:     "Cancel",
		Comments_ReplyingTo: "Replying to",

		Post_Views: "views",

		About_Intro:    "An open engineering notebook by a 5th-semester information systems student focusing on Go backend development, Python, Linux systems, and distributed architecture.",
		About_Building: "What I'm building",
		About_Stack:    "Stack",
		About_Contact:  "Contact",

		Contact_FormTitle:    "Send a message",
		Contact_NameLabel:    "Your name",
		Contact_EmailLabel:   "Your email",
		Contact_MessageLabel: "Message...",
		Contact_Submit:       "Send",
		Contact_Success:      "Message sent. I'll get back to you soon.",
		Contact_Error:        "Failed to send. Please try again or email me directly.",

		NotFound_Title: "Page not found",
		NotFound_Body:  "The page you're looking for doesn't exist or has been moved.",
		NotFound_Home:  "Go home",

		Forbidden_Title: "Access forbidden",
		Forbidden_Body:  "You don't have permission to access this page or resource.",
		Forbidden_Home:  "Go home",

		ServerError_Title: "Internal server error",
		ServerError_Body:  "Something went wrong on our end. Please try again later.",
		ServerError_Home:  "Go home",

		Error_SearchPlaceholder: "Search entire archive...",
		Error_Diagnostics:       "SYSTEM DIAGNOSTICS",

		Uses_Hardware: "Hardware",
		Uses_Software: "Software",
		Uses_Desk:     "Desk",

		Now_WorkingOn: "Working on",
		Now_Reading:   "Reading",
		Now_Learning:  "Learning",

		Reading_Title:    "Reading List",
		Reading_Subtitle: "Saved dispatches and reference material bookmarked for later.",

		Stats_PageTitle: "Site statistics and traffic breakdown.",

		Saved_Title:    "Saved Reading List",
		Saved_Subtitle: "Bookmarked articles stored locally on your device.",
		Saved_Empty:    "No saved posts yet. Click the bookmark icon on any post to save it for later.",
		Saved_Clear:    "Clear all",

		Lang_Switch: "ID",

		Search_Title:      "Search",
		Search_ResultsFor: "Results for",

		Guestbook_Title:     "Guestbook",
		Guestbook_NameLabel: "Name",
		Guestbook_BodyLabel: "Message (max 500 chars)...",
		Guestbook_Submit:    "Sign",
		Guestbook_Empty:     "No messages yet. Be the first to sign the guestbook!",

		Stats_Title: "Traffic statistics",

		Changelog_Title: "Changelog",
		Links_Title:     "Links",
		Links_Subtitle:  "Curated resources I find useful and return to often.",

		Terms_Title:       "Terms of Use",
		Terms_LastUpdated: "Last updated: August 2026",
	},
	"id": {
		Nav_Home:     "Beranda",
		Nav_Behind:   "behind",
		Nav_Terminal: "terminal",
		Footer_Copy:  "© 2026 Dafa",

		Behind_Title:        "Behind this website",
		Behind_Colophon:     "Kolofon",
		Behind_ColophonBody: "Situs ini satu binary Go: chi untuk routing, templ untuk template HTML, HTMX untuk bagian interaktif, Tailwind untuk styling, SQLite untuk komentar dan views. Server-side rendered, self-hosted, tanpa tracker.",
		Behind_Uses:         "Stack & uses",
		Behind_Learning:     "Now / sedang belajar",

		Footer_Bio:       "Buku catatan rekayasa dan portofolio terbuka. Mendalami Go backend, Python, sistem Linux, dan arsitektur backend.",
		Footer_Read:      "Navigasi",
		Footer_Site:      "Situs",
		Footer_Topics:    "Topik",
		Footer_Lab:       "Proyek & Lab",
		Footer_Resources: "Sumber Daya",
		Footer_Meta:      "Metadata",
		Footer_Legal:     "Lisensi & Kebijakan",
		Footer_Community: "Open Source",
		Footer_Ecosystem: "Feeds & Meta",
		Footer_Other:     "Lainnya & Tools",
		Privacy_Title:    "Kebijakan Privasi",

		Home_BehindDoor: "look behind this website",

		Blog_Title:             "Blog",
		Blog_ReadMore:          "Baca",
		Blog_MinRead:           "mnt baca",
		Blog_NoPosts:           "Belum ada post.",
		Blog_SearchPlaceholder: "Cari post...",
		Blog_NoResults:         "Tidak ada post yang cocok.",
		Blog_PrevPage:          "Lebih baru",
		Blog_NextPage:          "Lebih lama",
		Blog_LoadMore:          "Muat lebih banyak",

		Post_Tags:    "Tag",
		Post_Back:    "Kembali ke blog",
		Post_Related: "Post terkait",
		Post_TOC:     "Daftar Isi",

		Comments_Heading:    "Komentar",
		Comments_Empty:      "Belum ada komentar. Jadi yang pertama berkomentar!",
		Comments_NameLabel:  "Nama kamu",
		Comments_BodyLabel:  "Tulis komentar...",
		Comments_Submit:     "Kirim komentar",
		Comments_Error:      "Gagal mengirim komentar. Pastikan nama dan komentar terisi.",
		Comments_Delete:     "Hapus",
		Comments_Reply:      "Balas",
		Comments_Cancel:     "Batal",
		Comments_ReplyingTo: "Membalas",

		Post_Views: "dilihat",

		About_Intro:    "Buku catatan terbuka mahasiswa sistem informasi semester 5 yang mendalami Go backend, Python, sistem Linux, dan arsitektur backend. Dibuat untuk mendokumentasikan proses belajar dan eksperimen kode.",
		About_Building: "Yang sedang saya bangun",
		About_Stack:    "Stack",
		About_Contact:  "Kontak",

		Contact_FormTitle:    "Kirim pesan",
		Contact_NameLabel:    "Nama kamu",
		Contact_EmailLabel:   "Email kamu",
		Contact_MessageLabel: "Pesan...",
		Contact_Submit:       "Kirim",
		Contact_Success:      "Pesan terkirim! Saya akan segera membalas.",
		Contact_Error:        "Gagal mengirim. Coba email langsung ke realdaemontalk@gmail.com",

		NotFound_Title: "Halaman tidak ditemukan",
		NotFound_Body:  "Halaman yang kamu cari tidak ada atau sudah dipindahkan.",
		NotFound_Home:  "Ke beranda",

		Forbidden_Title: "Akses ditolak",
		Forbidden_Body:  "Kamu tidak memiliki izin untuk mengakses halaman ini.",
		Forbidden_Home:  "Ke beranda",

		ServerError_Title: "Kesalahan server",
		ServerError_Body:  "Terjadi kesalahan pada server. Silakan coba beberapa saat lagi.",
		ServerError_Home:  "Ke beranda",

		Error_SearchPlaceholder: "Cari di seluruh arsip...",
		Error_Diagnostics:       "DIAGNOSTIK SISTEM",

		Lang_Switch: "EN",

		Search_Title:      "Cari",
		Search_ResultsFor: "Hasil untuk",

		Guestbook_Title:     "Buku Tamu",
		Guestbook_NameLabel: "Nama kamu",
		Guestbook_BodyLabel: "Tinggalkan pesan...",
		Guestbook_Submit:    "Tanda tangani buku tamu",
		Guestbook_Empty:     "Belum ada pesan. Jadilah yang pertama!",

		Stats_Title: "Statistik",

		Changelog_Title: "Catatan Perubahan",
		Links_Title:     "Tautan",
		Links_Subtitle:  "Sumber daya kurasi yang sering saya kunjungi.",

		Terms_Title:       "Ketentuan Penggunaan",
		Terms_LastUpdated: "Terakhir diperbarui: Agustus 2026",
	},
}

func Get(lang string) UI {
	if s, ok := Strings[lang]; ok {
		return s
	}
	return Strings["en"]
}
