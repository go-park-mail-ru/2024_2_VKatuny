package servery

package main

import (
	"fmt"
	getsess "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/session/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"context"
)


type SurveyManager struct {
	Survey.UnimplementedAuthCheckerServer
}

func NewSuveyService() *SurveyManager {
	return &SurveyManager{
	}
}

func (sm *SurveyManager) GetAuthorization(in *GetAuthorizationInput) returns (*GetAuthorizationOutput) {
	authorized := &GetAuthorizationOutput{authorized: getsess.GetUserIdBySession(in.Token)}
	return authorized
}