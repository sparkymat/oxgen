package service

type DatabaseProvider interface{
  
  CountSearchedUsers(ctx context.Context, query string) ([]int64, error) 
  
  CreateUser(ctx context.Context, params dbx.CreateUserParams) (dbx.User, error)
  DeleteUser(ctx context.Context, id uuid.UUID) error 
  FetchUserByID(ctx context.Context, id uuid.UUID) (dbx.User, error) 
  FetchUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]dbx.User, error) 
  
  SearchUsers(ctx context.Context, arg dbx.SearchUsersParams) ([]dbx.User, error) 
  
  UpdateUserName(ctx context.Context, arg dbx.UpdateUserNameParams) (dbx.User, error)
  UpdateUserAge(ctx context.Context, arg dbx.UpdateUserAgeParams) (dbx.User, error)
  UpdateUserDob(ctx context.Context, arg dbx.UpdateUserDobParams) (dbx.User, error)
  UpdateUserPhotoPath(ctx context.Context, arg dbx.UpdateUserPhotoPathParams) (dbx.User, error)
}