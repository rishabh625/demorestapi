package object

//SignupRequest object for SignUp Request
type SignupRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	AdminSecret string `json:"admin_secret"`
	Email       string `json:"email"`
}

//Response object for Common Response
type Response struct {
	Message string `json:"message,omitempty"`
}

//LoginRequest object for Login Request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//LoginResponse object for Login Response
type LoginResponse struct {
	Message     string `json:"message,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

//AddMovieRequest Object For Add Movie Request
type AddMovieRequest struct {
	Movies []Movie `json:"movies"`
}

//Movie Movies data structure
type Movie struct {
	Id           int      `json:"id,omitempty"`
	Popularity99 float64  `json:"99popularity,omitempty"`
	Director     string   `json:"director,omitempty"`
	Genre        []string `json:"genre,omitempty"`
	ImdbScore    float64  `json:"imdb_score,omitempty"`
	Name         string   `json:"name,omitempty"`
}

//GetAllMovieResponse Object for All Movies List Response
type GetAllMovieResponse struct {
	Movies  []Movie `json:"movies"`
	Message string  `json:"message"`
}

//UpdateMovie Object for Update Movie Request
type UpdateMovie struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

//SearchMovie Object for Search Movie
type SearchMovie struct {
	Term  Movie   `json:"term,omitempty"`
	Range []Range `json:"range,omitempty"`
}

//Range Object For Range Query Structure
type Range struct {
	Key      string  `json:"key,omitempty"`
	Operator string  `json:"operator,omitempty"`
	Value    float64 `json:"value,omitempty"`
}
