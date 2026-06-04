package handler

import (
	"encoding/json"
	"net/http"

	"gitstat/internal/model"
	"gitstat/internal/scanner"
	"gitstat/internal/store"
)

// SetScanPathHandler 设置扫描路径，不立即扫描，等待懒加载
func SetScanPathHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	// 清空旧缓存
	store.GlobalStore.ClearAll()

	// 保存扫描路径
	store.GlobalStore.SetScanPath(req.Path)

	// 扫描新路径下的仓库元数据（不获取提交）
	repos, err := scanner.DiscoverRepos(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 注册仓库（未初始化状态）
	store.GlobalStore.RegisterRepos(repos)

	writeJSON(w, "SetScanPath", model.ApiResponse{
		Code: 200,
		Data: map[string]interface{}{
			"path": req.Path,
		},
		Message: "Path set successfully, data will be loaded on demand",
	})
}

func GetRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	repos := store.GlobalStore.GetRepositories()
	writeJSON(w, "Repositories", repos)
}

// GetScanPathHandler 获取当前扫描路径
func GetScanPathHandler(w http.ResponseWriter, r *http.Request) {
	path := store.GlobalStore.GetScanPath()
	writeJSON(w, "GetScanPath", model.ApiResponse{
		Code: 200,
		Data: map[string]string{
			"path": path,
		},
	})
}
