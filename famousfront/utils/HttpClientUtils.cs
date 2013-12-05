using famousfront.datamodels;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading.Tasks;

namespace famousfront.utils
{
  class ResultWithError<TDataType>
  {
    internal TDataType data = default(TDataType);
    internal int code = -1;
    internal string reason;
  }
  class HttpClientUtils
  {
    internal static async Task<ResultWithError<TDataType>> Get<TDataType>(string url)
    {
      using (var client = new HttpClient())
      {
        try
        {
          client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
          var resp = await client.GetAsync(url);
          var sc = resp.StatusCode;
          if (sc != System.Net.HttpStatusCode.OK)
          {
            var r = await resp.Content.ReadAsAsync<BackendError>();
            return new ResultWithError<TDataType> { code = r.code, reason = r.reason };
          }
          var r2 = await resp.Content.ReadAsAsync<TDataType>();
          return new ResultWithError<TDataType> { data = r2, code = 0, reason = resp.ReasonPhrase };
        }
        catch (Newtonsoft.Json.JsonException e)
        {
          return new ResultWithError<TDataType> { code = -1, reason = e.Message };
        }
        catch (HttpRequestException e)
        {
          return new ResultWithError<TDataType> { code = -2, reason = e.Message };
        }
        catch (System.Net.WebException e)
        {
          return new ResultWithError<TDataType> { code = -3, reason = e.Message };
        }
        catch (TaskCanceledException e)
        {
          return new ResultWithError<TDataType> { code = -4, reason = e.Message };
        }
      }

    }
  }

}
