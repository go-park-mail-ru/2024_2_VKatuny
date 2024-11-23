package servery

package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/go-park-mail-ru/lectures/8-microservices/4_grpc/Survey"

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

func (sm *SurveyManager) GetAuthorization(GetAuthorizationInput) returns (GetAuthorizationOutput) {
	authorized := GetAuthorizationOutput{}
	return authorized
}