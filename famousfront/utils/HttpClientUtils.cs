using famousfront.datamodels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Runtime.Serialization;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.utils
{
    class ResultWithError<DataType>
    {
        internal DataType data = default(DataType);
        internal int code = -1;
        internal string reason;
    }
    class HttpClientUtils
    {
        internal static async Task<ResultWithError<DataType>> Get<DataType>(string url)
        {
            try
            {
                using (var client = new HttpClient())
                {
                    client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
                    var resp = await client.GetAsync(url);
                    var sc = resp.StatusCode;
                    if (sc != System.Net.HttpStatusCode.OK)
                    {
                        var r = await resp.Content.ReadAsAsync<BackendError>();
                        return new ResultWithError<DataType> { code = r.code, reason = r.reason };
                    }
                    var r2 = await resp.Content.ReadAsAsync<DataType>();
                    return new ResultWithError<DataType> { data = r2, code = 0, reason = resp.ReasonPhrase };
                }
            }
            catch (Newtonsoft.Json.JsonException e)
            {
                return new ResultWithError<DataType> { code = -1, reason = e.Message };
            }
            catch (HttpRequestException e)
            {
                return new ResultWithError<DataType> { code = -2, reason = e.Message };
            }            
        }
    }
}
