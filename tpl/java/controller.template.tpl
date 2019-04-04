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
package {{.ProjectPkg}}.controller;

import com.airparking.cloud.common.AbstractController;
import com.airparking.cloud.common.dao.AbstractDAO;
import {{.ProjectPkg}}.{{.ModelPkg}}.{{.ModelName}};
import {{.ProjectPkg}}.service.{{.ModelName}}Service;

import com.alibaba.fastjson.JSONObject;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

/**
 * Created by {{.Author}} on {{.DateStr}}.
 */
@RestController
public interface {{.ModelName}}Controller extends AbstractController {
    @Autowired
    private {{.ModelName}}Service {{ToLowerCamel .ModelName}}Service;

    @RequestMapping("get")
    @ResponseBody
    public {{.ModelName}} get(@RequestParam("id") Long id) {
        return this.{{ToLowerCamel .ModelName}}Service.get(id);
    }

    @RequesetMapping("add")
    public Integer add(@RequestParam("body") String {{ToLowerCamel .ModelName}}Json) {
        return this.{{ToLowerCamel .ModelName}}Service.insert(JSONObject.parseObject({{ToLowerCamel .ModelName}}Json, Map.class));
    }
}