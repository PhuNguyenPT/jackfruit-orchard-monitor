package views

type Translations struct {
	// Nav
	Home       string
	Contact    string
	Products   string
	SignOut    string
	GetStarted string
	// Register
	RegisterTitle       string
	RegisterDescription string
	RegisterHeading     string
	NameLabel           string
	EmailLabel          string
	PasswordLabel       string
	RegisterButton      string
	AlreadyHaveAccount  string
	SignIn              string
	// Login
	LoginTitle       string
	LoginDescription string
	LoginHeading     string
	NoAccount        string
	Register         string
	LoginButton      string
	// Contact
	ContactTitle     string
	ContactHeading   string
	ContactSubtitle  string
	FullName         string
	Subject          string
	Message          string
	SendMessage      string
	ContactThanks    string
	ContactSent      string
	ContactFail      string
	ContactRateLimit string
	// Placeholders
	PlaceholderName     string
	PlaceholderEmail    string
	PlaceholderSubject  string
	PlaceholderMessage  string
	PlaceholderPassword string
	SearchPlaceholder   string
	// Dashboard
	DashboardTitle       string
	DashboardDescription string
	DashboardWelcome     string
	AccountInfo          string
	MemberSince          string
	ChangeName           string
	NewName              string
	UpdateName           string
	ChangePassword       string
	CurrentPassword      string
	NewPassword          string
	ConfirmPassword      string
	UpdatePassword       string
	Security             string
	ActiveSessions       string
	Revoke               string
	NameUpdated          string
	PasswordUpdated      string
	Show                 string
	Hide                 string
	Expires              string
	// Common
	GoHome string
	// Pagination
	Page string
	Of   string
	// Footer
	FooterDesc        string
	QuickLinks        string
	AllRightsReserved string
	FooterEmail       string
	FooterPhone       string
	// Error message
	ErrRequired         string
	ErrInvalidEmail     string
	ErrPasswordMismatch string
	ErrWrongPassword    string
	ErrEmailInUse       string
	ErrSomethingWrong   string
	ErrNameRequired     string
	ErrAllRequired      string
	ErrInvalidPassword  string
	// Home
	HomeWelcome        string
	HomeSubtitle       string
	HomeCTA            string
	HomeLatestProducts string
	HomeViewAll        string
	// Sensors
	SensorsTitle        string
	SensorsDescription  string
	SensorsHeading      string
	SensorsSubtitle     string
	SensorAddr          string
	SensorTemperature   string
	SensorHumidity      string
	SensorLastUpdated   string
	SensorNoData        string
	SensorAutoRefresh   string
	SHT40SectionHeading string
	SoilSectionHeading  string
	SoilMoistureLabel   string
	SensorRaw           string
	DeviceConnected     string
	DeviceDisconnected  string
	HumidityScaleLabel  string
	// Sensor History
	HistoryBack            string
	SHT40HistoryTitle      string
	SoilHistoryTitle       string
	TemperatureScaleLabel  string
	SoilMoistureScaleLabel string
	// About
	AboutNav               string
	AboutTitle             string
	AboutDescription       string
	AboutHeading           string
	AboutSubtitle          string
	AboutMission           string
	AboutMissionText       string
	AboutStack             string
	AboutHardware          string
	AboutWhy               string
	AboutWhyText           string
	AboutStackPostgresDesc string
	AboutStackHtmxDesc     string
	AboutStackMqttDesc     string
	AboutHardwareDesc      string
	AboutHardwareEsp32Desc string
	AboutHardwareSht40Desc string
	AboutHardwareSoilDesc  string
}

const (
	LangEN = "en"
	LangVI = "vi"
)

func T(lang string) Translations {
	if lang == LangVI {
		return Translations{
			// Nav
			Home:       "Trang chủ",
			Contact:    "Liên hệ",
			Products:   "Sản phẩm",
			SignOut:    "Đăng xuất",
			GetStarted: "Bắt đầu",
			// Register
			RegisterTitle:       "Đăng ký",
			RegisterDescription: "Tạo tài khoản mới",
			RegisterHeading:     "Tạo tài khoản",
			NameLabel:           "Họ tên",
			EmailLabel:          "Email",
			PasswordLabel:       "Mật khẩu",
			RegisterButton:      "Đăng ký",
			AlreadyHaveAccount:  "Đã có tài khoản?",
			SignIn:              "Đăng nhập",
			// Login
			LoginTitle:       "Đăng nhập",
			LoginDescription: "Đăng nhập vào tài khoản của bạn",
			LoginHeading:     "Đăng nhập",
			NoAccount:        "Chưa có tài khoản?",
			Register:         "Đăng ký",
			LoginButton:      "Đăng nhập",
			// Contact
			ContactTitle:     "Liên hệ",
			ContactHeading:   "Liên hệ với chúng tôi",
			ContactSubtitle:  "Bạn có câu hỏi? Chúng tôi rất vui được lắng nghe.",
			FullName:         "Họ tên",
			Subject:          "Tiêu đề",
			Message:          "Nội dung",
			SendMessage:      "Gửi tin nhắn",
			ContactThanks:    "Cảm ơn ",
			ContactSent:      "Tin nhắn của bạn đã được gửi.",
			ContactFail:      "Đã có lỗi xảy ra. Vui lòng thử lại sau.",
			ContactRateLimit: "Bạn đã đạt giới hạn tin nhắn trong ngày. Vui lòng thử lại vào ngày mai.",
			// Placeholders
			PlaceholderName:     "Nguyễn Văn A",
			PlaceholderEmail:    "a.nguyenvan@email.com",
			PlaceholderSubject:  "Chúng tôi có thể giúp gì cho bạn?",
			PlaceholderMessage:  "Nhập nội dung tin nhắn...",
			PlaceholderPassword: "••••••••",
			SearchPlaceholder:   "Tìm kiếm sản phẩm...",
			// Dashboard
			DashboardTitle:       "Bảng điều khiển",
			DashboardDescription: "Trang quản lý của bạn",
			DashboardWelcome:     "Xin chào",
			AccountInfo:          "Thông tin tài khoản",
			MemberSince:          "Thành viên từ",
			ChangeName:           "Đổi tên",
			NewName:              "Tên mới",
			UpdateName:           "Cập nhật tên",
			ChangePassword:       "Đổi mật khẩu",
			CurrentPassword:      "Mật khẩu hiện tại",
			NewPassword:          "Mật khẩu mới",
			ConfirmPassword:      "Xác nhận mật khẩu mới",
			UpdatePassword:       "Cập nhật mật khẩu",
			Security:             "Bảo mật",
			ActiveSessions:       "Phiên đang hoạt động",
			Revoke:               "Thu hồi",
			NameUpdated:          "Cập nhật tên thành công.",
			PasswordUpdated:      "Cập nhật mật khẩu thành công.",
			Show:                 "Hiện",
			Hide:                 "Ẩn",
			Expires:              "hết hạn",
			// Common
			GoHome: "Về trang chủ",
			// Pagination
			Page: "Trang",
			Of:   "của",
			// Footer
			FooterDesc:        "Hệ thống giám sát thông số môi trường thời gian thực qua IoT.",
			QuickLinks:        "Liên kết nhanh",
			AllRightsReserved: "Mọi quyền được bảo lưu.",
			FooterEmail:       "Email",
			FooterPhone:       "Điện thoại",
			// Error message
			ErrRequired:         "là bắt buộc.",
			ErrInvalidEmail:     "Vui lòng nhập địa chỉ email hợp lệ.",
			ErrPasswordMismatch: "Mật khẩu không khớp.",
			ErrWrongPassword:    "Mật khẩu hiện tại không đúng.",
			ErrEmailInUse:       "Email đã được sử dụng.",
			ErrSomethingWrong:   "Đã có lỗi xảy ra, vui lòng thử lại.",
			ErrNameRequired:     "Tên là bắt buộc.",
			ErrAllRequired:      "Vui lòng điền đầy đủ thông tin.",
			ErrInvalidPassword:  "Email hoặc mật khẩu không đúng.",
			// Home
			HomeWelcome:        "Giám Sát IoT",
			HomeSubtitle:       "Dữ liệu thời gian thực",
			HomeCTA:            "Xem bảng điều khiển",
			HomeLatestProducts: "Sản phẩm mới nhất",
			HomeViewAll:        "Xem tất cả →",
			// Sensors
			SensorsTitle:        "Cảm biến",
			SensorsDescription:  "Dữ liệu thời gian thực từ các cảm biến",
			SensorsHeading:      "Dữ liệu cảm biến",
			SensorsSubtitle:     "Tự động cập nhật mỗi 10 giây",
			SensorAddr:          "Cảm biến",
			SensorTemperature:   "Nhiệt độ",
			SensorHumidity:      "Độ ẩm",
			SensorLastUpdated:   "Cập nhật lúc",
			SensorNoData:        "Chưa có dữ liệu từ cảm biến.",
			SHT40SectionHeading: "SHT40 — Nhiệt độ & Độ ẩm không khí",
			SoilSectionHeading:  "MKE-S13 — Độ ẩm đất",
			SoilMoistureLabel:   "Độ ẩm đất",
			HumidityScaleLabel:  "Độ ẩm không khí",
			SensorRaw:           "Giá trị gốc",
			DeviceConnected:     "Đã kết nối",
			DeviceDisconnected:  "Mất kết nối",
			// Sensor History
			SensorAutoRefresh:      "Tự động cập nhật",
			HistoryBack:            "Quay lại",
			SHT40HistoryTitle:      "SHT40 — Cảm biến",
			SoilHistoryTitle:       "MKE-S13 — Cảm biến",
			TemperatureScaleLabel:  "Thang nhiệt độ",
			SoilMoistureScaleLabel: "Thang độ ẩm đất",
			// About
			AboutNav:               "Giới thiệu",
			AboutTitle:             "Giới thiệu — Prizm",
			AboutDescription:       "Giám sát IoT thời gian thực bằng Go, WebSocket, PostgreSQL, HTMX và Templ.",
			AboutHeading:           "Về dự án Prizm",
			AboutSubtitle:          "Nền tảng IoT kết nối phần cứng thực tế với kiến trúc Go hiện đại.",
			AboutMission:           "Sứ mệnh",
			AboutMissionText:       "Prizm cung cấp dữ liệu môi trường thời gian thực từ các cảm biến IoT triển khai ngoài thực địa — bắt đầu từ vườn mít ở Việt Nam.",
			AboutStack:             "Công nghệ",
			AboutHardware:          "Phần cứng",
			AboutWhy:               "Mục tiêu ra đời",
			AboutWhyText:           "Prizm ra đời như một dự án thực chiến nhằm ứng dụng lập trình Go Full-stack vào thiết bị thực tế — cầu nối giữa lập trình nhúng (Firmware ESP32) và kiến trúc Backend hiện đại.",
			AboutStackPostgresDesc: "Truy vấn Type-safe với sqlc",
			AboutStackHtmxDesc:     "SSR & Giao diện tương tác tức thì",
			AboutStackMqttDesc:     "Truyền dữ liệu thời gian thực qua WebSocket",
			AboutHardwareDesc:      "Sự kết hợp giữa vi điều khiển hiệu năng cao và các cảm biến tiêu chuẩn công nghiệp giúp theo dõi chính xác môi trường nhà vườn theo thời gian thực, đồng thời gửi dữ liệu trực tiếp về hệ thống Backend Go.",
			AboutHardwareEsp32Desc: "Bộ vi điều khiển trung tâm tích hợp Wi-Fi, chịu trách nhiệm tổng hợp dữ liệu và truyền tải an toàn qua giao thức MQTT/TLS.",
			AboutHardwareSht40Desc: "Cảm biến đo nhiệt độ và độ ẩm với độ chính xác cao, sử dụng chuẩn giao tiếp Modbus RTU.",
			AboutHardwareSoilDesc:  "Cảm biến điện dung đo độ ẩm đất, được mở rộng kênh đọc tín hiệu thông qua IC 74HC4067.",
		}
	}
	return Translations{
		// Nav
		Home:       "Home",
		Contact:    "Contact",
		Products:   "Products",
		SignOut:    "Sign out",
		GetStarted: "Get Started",
		// Register
		RegisterTitle:       "Register",
		RegisterDescription: "Create a new account",
		RegisterHeading:     "Create an Account",
		NameLabel:           "Name",
		EmailLabel:          "Email",
		PasswordLabel:       "Password",
		RegisterButton:      "Register",
		AlreadyHaveAccount:  "Already have an account?",
		SignIn:              "Sign in",
		// Login
		LoginTitle:       "Login",
		LoginDescription: "Sign in to your account",
		LoginHeading:     "Sign In",
		NoAccount:        "Don't have an account?",
		Register:         "Register",
		LoginButton:      "Sign In",
		// Contact
		ContactTitle:     "Contact Us",
		ContactHeading:   "Get in Touch",
		ContactSubtitle:  "Have a question? We'd love to hear from you.",
		FullName:         "Full Name",
		Subject:          "Subject",
		Message:          "Message",
		SendMessage:      "Send Message",
		ContactThanks:    "Thank you ",
		ContactSent:      "Your message has been sent.",
		ContactFail:      "Something went wrong. Please try again later.",
		ContactRateLimit: "You've reached the maximum number of messages for today. Please try again tomorrow.",
		// Placeholders
		PlaceholderName:     "John Doe",
		PlaceholderEmail:    "john@example.com",
		PlaceholderSubject:  "How can we help you?",
		PlaceholderMessage:  "Write your message here...",
		PlaceholderPassword: "••••••••",
		SearchPlaceholder:   "Search products...",
		// Dashboard
		DashboardTitle:       "Dashboard",
		DashboardDescription: "Your dashboard",
		DashboardWelcome:     "Welcome",
		AccountInfo:          "Account Info",
		MemberSince:          "Member since",
		ChangeName:           "Change Name",
		NewName:              "New Name",
		UpdateName:           "Update Name",
		ChangePassword:       "Change Password",
		CurrentPassword:      "Current Password",
		NewPassword:          "New Password",
		ConfirmPassword:      "Confirm New Password",
		UpdatePassword:       "Update Password",
		Security:             "Security",
		ActiveSessions:       "Active Sessions",
		Revoke:               "Revoke",
		NameUpdated:          "Name updated successfully.",
		PasswordUpdated:      "Password updated successfully.",
		Show:                 "Show",
		Hide:                 "Hide",
		Expires:              "expires",
		// Common
		GoHome: "Go to Home",
		// Pagination
		Page: "Page",
		Of:   "of",
		// Footer
		FooterDesc:        "Real-time monitoring powered by IoT sensors.",
		QuickLinks:        "Quick Links",
		AllRightsReserved: "All rights reserved.",
		FooterEmail:       "Email",
		FooterPhone:       "Phone",
		// Error message
		ErrRequired:         "is required.",
		ErrInvalidEmail:     "Please enter a valid email address.",
		ErrPasswordMismatch: "Passwords do not match.",
		ErrWrongPassword:    "Current password is incorrect.",
		ErrEmailInUse:       "Email already in use.",
		ErrSomethingWrong:   "Something went wrong, please try again.",
		ErrNameRequired:     "Name is required.",
		ErrAllRequired:      "All fields are required.",
		ErrInvalidPassword:  "Invalid email or password.",
		// Home
		HomeWelcome:        "IoT Monitor",
		HomeSubtitle:       "Real-time data monitoring",
		HomeCTA:            "View Dashboard",
		HomeLatestProducts: "Latest Products",
		HomeViewAll:        "View all →",
		// Sensors
		SensorsTitle:        "Sensors",
		SensorsDescription:  "Live readings from all sensors",
		SensorsHeading:      "Sensor Readings",
		SensorsSubtitle:     "Auto-refreshes every 10 seconds",
		SensorAddr:          "Sensor",
		SensorTemperature:   "Temperature",
		SensorHumidity:      "Humidity",
		SensorLastUpdated:   "Updated",
		SensorNoData:        "No sensor data available yet.",
		SensorAutoRefresh:   "Auto-refresh",
		SHT40SectionHeading: "SHT40 — Temperature & Air Humidity",
		SoilSectionHeading:  "MKE-S13 — Soil Moisture",
		SoilMoistureLabel:   "Soil Moisture",
		HumidityScaleLabel:  "Air Humidity",
		SensorRaw:           "Raw",
		DeviceConnected:     "Connected",
		DeviceDisconnected:  "Disconnected",
		// Sensor History
		HistoryBack:            "Back",
		SHT40HistoryTitle:      "SHT40 — Sensor",
		SoilHistoryTitle:       "MKE-S13 — Sensor",
		TemperatureScaleLabel:  "Temperature Scale",
		SoilMoistureScaleLabel: "Soil Moisture Scale",
		// About
		AboutNav:               "About",
		AboutTitle:             "About — Prizm",
		AboutDescription:       "Real-time IoT monitoring built with Go, WebSocket, PostgreSQL, HTMX and Templ.",
		AboutHeading:           "About Prizm",
		AboutSubtitle:          "Real-time IoT monitoring platform",
		AboutMission:           "Mission",
		AboutMissionText:       "Prizm delivers real-time environmental data from IoT sensors deployed in the field — starting with a jackfruit orchard in Vietnam.",
		AboutStack:             "Tech Stack",
		AboutHardware:          "Hardware",
		AboutWhy:               "Why We Built This",
		AboutWhyText:           "Prizm started as an educational IoT project to learn full-stack Go with real hardware — combining ESP32 embedded firmware with a modern Go web backend.",
		AboutStackPostgresDesc: "sqlc generated",
		AboutStackHtmxDesc:     "SSR + Dynamic UI",
		AboutStackMqttDesc:     "WebSocket real-time",
		AboutHardwareDesc:      "Industrial-grade sensors paired with robust microcontrollers ensure accurate, real-time monitoring of the orchard's environment, directly feeding into our Go backend.",
		AboutHardwareEsp32Desc: "The core WiFi microcontroller handling data aggregation and secure MQTT transmission over TLS.",
		AboutHardwareSht40Desc: "High-precision temperature and humidity tracking communicated via Modbus RTU protocol.",
		AboutHardwareSoilDesc:  "Capacitive soil moisture monitoring multiplexed securely via a 74HC4067 component.",
	}
}
