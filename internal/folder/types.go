package folder

import "github.com/google/uuid"

type CreateFolderRequest struct {
	Name string `json:"name" validate:"required,min=1,max=64"`
	Icon string `json:"icon"`
}

type UpdateFolderRequest struct {
	Name string `json:"name" validate:"required,min=1,max=64"`
	Icon string `json:"icon"`
}

type ReorderRequest struct {
	FolderIDs []uuid.UUID `json:"folderIds" validate:"required"`
}

type AddItemsRequest struct {
	ConversationIDs []uuid.UUID `json:"conversationIds" validate:"required,min=1"`
}
