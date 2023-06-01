/*Copyright [2023] [Alejandro Escanero Blanco <aescanero@disasterproject.com>]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.*/

/*

ports:
  ldap: 1389
  ldaps: 1686
  ca.pem: ...
  crt.pem: ...
  crt.key: ...
  ca.pem.file: ...
  crt.pem.file: ...
  crt.key.file: ...
databases:
- base: ...
  replicas:
  - url:
    ca.pem: ...
    crt.pem: ...
    crt.key: ...
	ca.pem.file: ...
    crt.pem.file: ...
    crt.key.file: ...
	attrs: "*,+"


*/
/*
olcSynrepl:
syncrepl
 ...
 provider=ldaps://ldap.example.com
 bindmethod=simple
 binddn="cn=goodguy,dc=example,dc=com"
 credentials=dirtysecret
 starttls=critical
 schemachecking=on
 scope=sub
 searchbase="dc=example,dc=com"
 tls_cacert=/path/to/file
 tls_cert=/path/to/file.ext
 tls_key=/path/to/file.ext
 tls_protocol_min=1.2
 tls_reqcert=demand
 type=refreshAndPersist

*/
package config

type Config struct {
	SrvConfig ServerConfig
	Database  []DatabaseConfig
	Schemas   []SchemaConfig
	Modules   []ModuleConfig
}
