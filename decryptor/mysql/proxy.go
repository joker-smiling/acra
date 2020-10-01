/*
Copyright 2018, Cossack Labs Limited

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysql

import (
	"github.com/cossacklabs/acra/decryptor/base"
	"github.com/cossacklabs/acra/encryptor"
)

type proxyFactory struct {
	dataEncryptor encryptor.DataEncryptor
	setting       base.ProxySetting
}

// NewProxyFactory return new proxyFactory
func NewProxyFactory(proxySetting base.ProxySetting) (base.ProxyFactory, error) {
	dataEncryptor, err := encryptor.NewAcrawriterDataEncryptor(proxySetting.KeyStore())
	if err != nil {
		return nil, err
	}
	return &proxyFactory{
		dataEncryptor: dataEncryptor,
		setting:       proxySetting,
	}, nil
}

// New return mysql proxy implementation
func (factory *proxyFactory) New(clientID []byte, clientSession base.ClientSession) (base.Proxy, error) {
	decryptor, err := factory.setting.DecryptorFactory().New(clientID)
	if err != nil {
		return nil, err
	}
	proxy, err := NewMysqlProxy(clientSession, decryptor, factory.setting)
	if err != nil {
		return nil, err
	}
	if !factory.setting.TableSchemaStore().IsEmpty() {
		queryEncryptor, err := encryptor.NewMysqlQueryEncryptor(factory.setting.TableSchemaStore(), clientID, factory.dataEncryptor)
		if err != nil {
			return nil, err
		}
		proxy.AddQueryObserver(queryEncryptor)
	}
	proxy.SubscribeOnAllColumnsDecryption(decryptor)
	return proxy, nil
}
