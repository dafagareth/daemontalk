package i18n

type UI struct {
	Nav_Home     string
	Nav_Colophon string
	Nav_Terminal string
	Footer_Copy  string

	Colophon_Title    string
	Colophon_Nav      string
	Colophon_Body     string
	Colophon_Uses     string
	Colophon_Learning string

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

	Home_ColophonDoor string

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

	Saved_Title        string
	Saved_Subtitle     string
	Saved_Empty        string
	Saved_Clear        string
	Saved_Count        string
	Saved_EmptyHeading string
	Saved_EmptyBody    string
	Saved_Remove       string

	Daily_Title      string
	Daily_Dispatches string

	Lang_Switch string

	Search_Title      string
	Search_ResultsFor string
	Search_EmptyText  string

	Stats_Title          string
	Stats_SecPub         string
	Stats_SecLocales     string
	Stats_SecCommunity   string
	Stats_SecEngine      string
	Stats_TotalArticles  string
	Stats_TotalWords     string
	Stats_SystemTags     string
	Stats_PageViews      string
	Stats_Comments       string
	Stats_Topics         string
	Stats_Replies        string
	Stats_Solved         string
	Stats_Members        string
	Stats_Votes          string
	Stats_LangID         string
	Stats_LangEN         string
	Stats_LangES         string
	Stats_DBStorage      string
	Stats_ServerArch     string
	Stats_TemplateEngine string

	Changelog_Title string
	Links_Title     string
	Links_Subtitle  string

	Terms_Title       string
	Terms_LastUpdated string

	Link_About         string
	Link_Colophon      string
	Link_Daily         string
	Link_Graph         string
	Link_Terminal      string
	Link_ReadingList   string
	Link_ContribGuide  string
	Link_Discussions   string
	Link_Shortcuts     string
	Link_Contribute    string
	Link_License       string
	Link_Accessibility string
	Link_Resume        string
	Link_Contact       string
	Discussions_Title  string
	Discussions_New    string

	Modal_ShortcutsTitle string
	Modal_Cancel         string
	Modal_Delete         string
	Search_SeeAll        string

	Nav_All          string
	Nav_LiveDispatch string
	Nav_Radar        string

	Portal_TopViews   string
	Portal_NoPopPosts string
	Portal_Curated    string
	Portal_SeeMore    string
	Portal_ReadMore   string

	Sidebar_FullReader string
	Sidebar_Themes     string
	Sidebar_WebLink    string
	Sidebar_RSS        string
	Sidebar_JSON       string
	Sidebar_AuthorBio  string

	Post_Current string
	Post_Serif   string

	Auth_ExportData    string
	Auth_DeleteAccount string
	Auth_DeleteConfirm string
}

func Get(lang string) UI {
	if s, ok := Strings[lang]; ok {
		return s
	}
	return Strings["en"]
}
