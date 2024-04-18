package states

type State struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Uri        string   `json:"uri"`
	Region     string   `json:"region"`
	Capital    string   `json:"capital"`
	Deputy     string   `json:"deputy"`
	Governor   string   `json:"governor"`
	Lgas       []string `json:"lgas"`
	Neighbours []string `json:"neighbours"`
	Slogan     string   `json:"slogan"`
	Towns      []string `json:"towns"`
}
