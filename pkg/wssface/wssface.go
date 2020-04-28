package wssface

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/face"
	"github.com/Azure/go-autorest/autorest"
)

//FaceServiceClient client for the Azure Face Service
type FaceServiceClient struct {
	client *face.Client
}

//NewFaceServiceClient FaceServiceClient constructor
func NewFaceServiceClient(conf Configuration) *FaceServiceClient {
	// Client used for Detect Faces, Find Similar, and Verify examples.
	client := face.NewClient(conf.FaceEndpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(conf.FaceSubscription)

	return &FaceServiceClient{
		client: &client,
	}
}

//InvokeFace invokes the Face APIs with the provided photo
func (s *FaceServiceClient) InvokeFace(faceContext context.Context, photo io.ReadCloser) (*FaceResult, error) {
	detectSingleFaces, err := s.callFaceService(faceContext, photo)
	if err != nil {
		return nil, err
	}

	dFaces := *detectSingleFaces.Value
	if numFaces := len(dFaces); numFaces == 0 {
		return nil, fmt.Errorf("no faces identified in provided photo")
	} else if numFaces > 1 {
		return nil, fmt.Errorf("more than one face (more precisely %v) identified in provided photo", numFaces)
	}

	dFace := dFaces[0]
	emotion := getEmotion(dFace.FaceAttributes)

	return &FaceResult{
		Age:       *dFace.FaceAttributes.Age,
		Gender:    fmt.Sprintf("%s", dFace.FaceAttributes.Gender),
		Sentiment: emotion,
	}, nil
}

func (s *FaceServiceClient) callFaceService(faceContext context.Context, photo io.ReadCloser) (face.ListDetectedFace, error) {
	// Use recognition model 2 for feature extraction. Recognition model 1 is used to simply recogize faces.
	recognitionModel02, detectionModel := face.Recognition02, face.Detection01
	// Array types chosen for the attributes of Face
	attributes := []face.AttributeType{"age", "emotion", "gender"}
	returnFaceID, returnRecognitionModel, returnFaceLandmarks := true, false, false

	return s.client.DetectWithStream(
		faceContext, photo, &returnFaceID, &returnFaceLandmarks,
		attributes, recognitionModel02, &returnRecognitionModel, detectionModel)
}

// FaceResult result of the FaceAPI
type FaceResult struct {
	Age       float64
	Sentiment Emotion
	Gender    string
}
