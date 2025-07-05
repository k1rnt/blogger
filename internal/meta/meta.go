package meta

type FrontMatter struct {
    Title     string   `yaml:"title"`
    Labels    []string `yaml:"labels"`
    BloggerID string   `yaml:"blogger_id,omitempty"`
    Slug      string   `yaml:"slug,omitempty"`
    Status    string   `yaml:"status,omitempty"`  // draft/publish
    Published string   `yaml:"published,omitempty"`
}
