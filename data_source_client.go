package main

import (
    "log"
    
    "net/http"
    "io/ioutil"
    "github.com/hashicorp/terraform/helper/schema"
    "encoding/json"
    "errors"
)

func dataSourceClient() *schema.Resource {
    return &schema.Resource {
        Read: dataSourceClientRead,

        Schema: map[string]*schema.Schema {
            "name": &schema.Schema {
                Type: schema.TypeString,
                Required: true,
            },
            "client_id": &schema.Schema {
                Type: schema.TypeString,
                Computed: true,
            },
        },
    }
}

type ClientRecord struct {
    Tenant string
    Name string
    Client_id string
}

func dataSourceClientRead(d *schema.ResourceData, meta interface{}) error {
    name := d.Get("name").(string)
    client := http.Client{}

    config := meta.(Config)
    req, _ := http.NewRequest("GET", "https://" + config.domain + "/api/v2/clients?fields=name%2Cclient_id", nil)
    req.Header.Add("Authorization", "Bearer " + config.accessToken)

    res, err := client.Do(req)

    if nil != err {
        return err
    }
    defer res.Body.Close()

    data, err := ioutil.ReadAll(res.Body)
    if nil != err {
        return err
    }

    log.Printf("[DEBUG] Response was: " + string(data))

    var clientRecords []ClientRecord
    err = json.Unmarshal(data, &clientRecords)
    if nil != err {
        return err
    }

    for _, clientRecord := range clientRecords {
        if clientRecord.Name == name {
            d.Set("client_id", clientRecord.Client_id)
            d.SetId(clientRecord.Client_id)
            return nil
        }
    }

    return errors.New("Error: client with name '" + name + "' not found in domain '" + config.domain + "'")
}

