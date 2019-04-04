/*
 * Copyright © 2019 Airparking HERE <ryan.cao@airparking.cn>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package {{.ProjectPkg}}.dao.impl;

import com.airparking.cloud.common.dao.impl.AbstractDAOImpl;
import {{.ProjectPkg}}.{{.ModelPkg}}.{{.ModelName}};
import {{.ProjectPkg}}.dao.{{.ModelName}}DAO;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;

/**
 * Created by {{.Author}} on {{.DateStr}}.
 */
@Repository
public interface {{.ModelName}}DAOImpl extends AbstractDAOImpl<{{.ModelName}}, Long> implements {{.ModelName}}DAO {
    @AutoWired
    private {{.ModelName}}Mapper {{ToLowerCamel .ModelName}}Mapper;

    @Override
    public Mapper<{{.ModelName}}> getMapper() {
        return {{ToLowerCamel .ModelName}}DAO;
    }
}