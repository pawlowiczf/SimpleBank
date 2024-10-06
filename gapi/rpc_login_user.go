package gapi

import (
	"context"
	"database/sql"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.GetUsername(), server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.GetUsername(), server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.GetUsername(),
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	response := &pb.LoginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
	}

	return response, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	//
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
