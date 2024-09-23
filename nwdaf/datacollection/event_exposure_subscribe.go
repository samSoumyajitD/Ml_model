// package datacollection

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"

// 	"github.com/ciromacedo/nwdaf/consumer"
// 	nwdaf_context "github.com/ciromacedo/nwdaf/context"
// 	"github.com/ciromacedo/nwdaf/util"
// 	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
// 	"github.com/free5gc/openapi/models"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func InitEventExposureSubscriber(self *nwdaf_context.NWDAFContext) {

// 	searchOpt := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}
// 	// recupera todas as AMFs registradas na NRF
// 	resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_AMF, models.NfType_NWDAF, searchOpt)
// 	//resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_NEF, models.NfType_NWDAF, searchOpt)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	//para cada uma das AMF's registrar no core realiza o subscriber de coleta
// 	for _, nfProfile := range resp.NfInstances {

// 		/* localiza a URL do end-point de subscriber com status de REGISTRADO */
// 		amfUri, endpoint, apiversion := util.SearchNFServiceUri(nfProfile, models.ServiceName_NAMF_EVTS, models.NfServiceStatus_REGISTERED)
// 		//nefUri, endpoint, apiversion := util.SearchNFServiceUri(nfProfile, models.ServiceName_NNEF_EVENTSSUBSCRIPTION, models.NfServiceStatus_REGISTERED)

// 		fmt.Println(endpoint)
// 		fmt.Println(apiversion)

// 		var buffer bytes.Buffer

// 		buffer.WriteString(amfUri)
// 		//buffer.WriteString(nefUri)
// 		buffer.WriteString("/")
// 		buffer.WriteString(endpoint)
// 		buffer.WriteString("/")
// 		buffer.WriteString(apiversion)
// 		buffer.WriteString("/")
// 		buffer.WriteString("subscriptions")

// 		url := buffer.String()

// 		/*
// 		 * 1 º os possiveis tipos de eventos p/ AMF estão em AmfEventType
// 		 */

// 		jsonData := `
//     {
// 		"Subscription" : { 	"EventList"	:
// 										[{ "Type" : "REGISTRATION_ACCEPT",
//                                            "ImmediateFlag" : true}],
// 							"EventNotifyUri": "http://127.0.0.56:29599/datacollection/amf-contexts/registration-accept",
// 							"AnyUE" : true,
// 							"NfId"  : "NWDAF"
//                           },
// 		"SupportedFeatures"	: "xx"
// 	}

// 	`

// 		var jsonStr = []byte(jsonData)
// 		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
// 		req.Header.Set("X-Custom-Header", "myvalue")
// 		req.Header.Set("Content-Type", "application/json")

// 		client := util.GetHttpConnection()

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		defer resp.Body.Close()

// 		body, err := ioutil.ReadAll(resp.Body) // response body is []byte
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Println("string starting")
// 		fmt.Println(string(body))
// 		fmt.Println("string end")

// 		// Connect to MongoDB
// 		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 		mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer mongoClient.Disconnect(context.TODO())

// 		// Ensure the MongoDB connection is established
// 		err = mongoClient.Ping(context.TODO(), nil)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// Get the userData collection
// 		collection := mongoClient.Database("free5gc").Collection("dataCollectionNwdar")

// 		// Insert the response body into the userData collection
// 		insertResult, err := collection.InsertOne(context.TODO(), bson.M{"response": string(body)})
// 		if err != nil {
// 			log.Fatal(err)
// 		}

//			fmt.Println("Inserted document ID: ", insertResult.InsertedID)
//		}
//	}
package datacollection

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ciromacedo/nwdaf/consumer"
	nwdaf_context "github.com/ciromacedo/nwdaf/context"
	"github.com/ciromacedo/nwdaf/util"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitEventExposureSubscriber(self *nwdaf_context.NWDAFContext) {

	searchOpt := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}
	// recupera todas as AMFs registradas na NRF
	resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_AMF, models.NfType_NWDAF, searchOpt)
	//resp, err := consumer.SendSearchNFInstances(self.NrfUri, models.NfType_NEF, models.NfType_NWDAF, searchOpt)
	if err != nil {
		fmt.Println(err)
	}

	//para cada uma das AMF's registrar no core realiza o subscriber de coleta
	for _, nfProfile := range resp.NfInstances {

		/* localiza a URL do end-point de subscriber com status de REGISTRADO */
		amfUri, endpoint, apiversion := util.SearchNFServiceUri(nfProfile, models.ServiceName_NAMF_EVTS, models.NfServiceStatus_REGISTERED)
		//nefUri, endpoint, apiversion := util.SearchNFServiceUri(nfProfile, models.ServiceName_NNEF_EVENTSSUBSCRIPTION, models.NfServiceStatus_REGISTERED)

		fmt.Println(endpoint)
		fmt.Println(apiversion)

		var buffer bytes.Buffer

		buffer.WriteString(amfUri)
		//buffer.WriteString(nefUri)
		buffer.WriteString("/")
		buffer.WriteString(endpoint)
		buffer.WriteString("/")
		buffer.WriteString(apiversion)
		buffer.WriteString("/")
		buffer.WriteString("subscriptions")

		url := buffer.String()

		/*
		 * 1 º os possiveis tipos de eventos p/ AMF estão em AmfEventType
		 */

		jsonData := `
    {	
		"Subscription" : { 	"EventList"	: 
										[{ "Type" : "REGISTRATION_ACCEPT",
                                           "ImmediateFlag" : true}], 
							"EventNotifyUri": "http://127.0.0.56:29599/datacollection/amf-contexts/registration-accept",
							"AnyUE" : true,
							"NfId"  : "NWDAF"
                          },
		"SupportedFeatures"	: "xx"
	}
		
	`

		var jsonStr = []byte(jsonData)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := util.GetHttpConnection()

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body) // response body is []byte
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("string starting")
		fmt.Println(string(body))
		fmt.Println("string end")

		// Connect to MongoDB
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer mongoClient.Disconnect(context.TODO())

		// Ensure the MongoDB connection is established
		err = mongoClient.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		// Get the dataCollectionNwdar collection
		collection := mongoClient.Database("free5gc").Collection("dataCollectionNwdar")

		// Update the document if it exists, otherwise insert a new document
		filter := bson.M{"nfId": "NWDAF"}
		update := bson.M{
			"$set": bson.M{"response": string(body)},
		}
		upsert := true
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update, &opts)
		if err != nil {
			log.Fatal(err)
		}

		if updateResult.MatchedCount > 0 {
			fmt.Println("Matched and updated an existing document")
		} else if updateResult.UpsertedCount > 0 {
			fmt.Printf("Inserted a new document with ID: %v\n", updateResult.UpsertedID)
		} else {
			fmt.Println("No documents matched, and no new documents were inserted")
		}
	}
}
