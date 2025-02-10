package domain

// PageAnalysis represents the result of webpage analysis
type PageAnalysis struct {
    HTMLVersion   string        `json:"htmlVersion"`
    PageTitle     string        `json:"pageTitle"`
    Headings      HeadingCount  `json:"headings"`
    Links         LinkAnalysis  `json:"links"`
    HasLoginForm  bool          `json:"hasLoginForm"`
}

// HeadingCount stores the count of different heading levels
type HeadingCount struct {
    H1 int `json:"h1"`
    H2 int `json:"h2"`
    H3 int `json:"h3"`
    H4 int `json:"h4"`
    H5 int `json:"h5"`
    H6 int `json:"h6"`
}

// LinkAnalysis represents the analysis of links in the webpage
type LinkAnalysis struct {
    Internal     int `json:"internal"`
    External     int `json:"external"`
    Inaccessible int `json:"inaccessible"`
}

// AnalysisRequest represents the incoming request for webpage analysis
type AnalysisRequest struct {
    URL string `json:"url" binding:"required,url"`
}