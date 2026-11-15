package tytd_go

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"tytd_go/internal/model"
)

// webFleet is a privacy-preserving account overview. Only the configured
// administrator can inspect all accounts. Other users see only their own
// account and their explicitly owned leader members.
func (t *TeleBot) webFleet(w http.ResponseWriter, r *http.Request, chatID int64, u *model.User) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	admin := t.isAdmin(chatID)
	db := t.ydb.GetDB()
	var accounts []model.User
	q := db.Model(&model.User{})
	if !admin {
		q = q.Where("id = ? OR leader_id = ?", strconv.FormatInt(chatID, 10), strconv.FormatInt(chatID, 10))
	}
	if err := q.Order("id asc").Find(&accounts).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprint(w, `<!doctype html><html lang="vi"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>TYTD Fleet</title><style>body{font:16px system-ui;background:#101827;color:#eef2ff;margin:0;padding:22px}main{max-width:950px;margin:auto}a{color:#93c5fd}table{width:100%;border-collapse:collapse;background:#1e293b}th,td{padding:12px;text-align:left;border-bottom:1px solid #334155} .scroll{overflow:auto} .pill{padding:3px 9px;border-radius:20px;background:#334155} </style></head><body><main><h1>TYTD · Quản lý tài khoản</h1><p><a href="/">Dashboard</a> · <a href="/members">Quản lý đệ tử</a> · <a href="/logout">Thoát</a></p><p>Danh sách: `)
	fmt.Fprintf(w, "%d", len(accounts))
	fmt.Fprint(w, ` tài khoản. Thao tác đăng nhập, cấu hình và bật/tắt auto được thực hiện trong trang Đệ tử của chủ sở hữu.</p><div class="scroll"><table><thead><tr><th>Mã</th><th>Tên game</th><th>Server</th><th>Chủ sở hữu</th><th>Auto</th><th>Trạng thái</th></tr></thead><tbody>`)
	for _, a := range accounts {
		owner := a.ID
		if a.LeaderID != "" {
			owner = a.LeaderID
		}
		name := a.IngameName
		if strings.TrimSpace(name) == "" {
			name = a.Username
		}
		running := "Đã dừng"
		if a.Running {
			running = "Đang chạy"
		}
		auto := "Tắt"
		if a.AutoLogin {
			auto = "Bật"
		}
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td><td><span class=\"pill\">%s</span></td><td>%s</td></tr>", html.EscapeString(a.ID), html.EscapeString(name), a.ServerID, html.EscapeString(owner), auto, running)
	}
	fmt.Fprint(w, `</tbody></table></div><p>Không hiển thị mật khẩu hoặc token. Chỉ admin được xem toàn hệ thống.</p></main></body></html>`)
}
