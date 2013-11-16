using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Text;
using famous.oauth.core;
using famous.oauth.core.request;

namespace famous.oauth.utils
{
  /// <summary>Utility class for building a URI using <see cref="BuildUri"/> or a HTTP request using 
  /// <see cref="BuildUri"/> from the query and path parameters of a REST call.</summary>
  public class RequestBuilder
  {
    /// <summary>
    /// A dictionary containing the parameters which will be inserted into the path
    /// of the URI. These parameters will be substituted into the URI path where the 
    /// path contains "{key}" that portion of the path will be replaced by the value 
    /// for the specified key in this dictionary.
    /// </summary>
    private IDictionary<string, string> PathParameters { get; set; }

    /// <summary>
    /// A dictionary containing the parameters which will apply to the query portion
    /// of this request.
    /// </summary>
    private List<KeyValuePair<string, string>> QueryParameters { get; set; }

    /// <summary>
    /// The base uri for this request (usually applies to the service itself).
    /// </summary>
    public Uri BaseUri { get; set; }

    /// <summary>
    /// The path portion of this request. Appended to the <see cref="BaseUri"/> and
    /// the parameters are substituted from the <see cref="PathParameters"/> dictionary.
    /// like /{id}/{uri}/{value}
    /// </summary>
    public string PathPattern { get; set; }

    /// <summary>Construct a new request builder.</summary> 
    public RequestBuilder(object request)
    {
      this.PathParameters = new Dictionary<string, string>();
      this.QueryParameters = new List<KeyValuePair<string, string>>();
      ParameterUtils.IterateParameters(request, (type, name, v) => AddParameter(type, name, v.ToString()));
    }

    /// <summary>Constructs a Uri as defined by the parts of this request builder.</summary> 
    public Uri BuildUri()
    {
      var restPath = new StringBuilder(PathParameters
        .Select(param => new { Token = "{" + param.Key + "}", Value = Uri.EscapeDataString(param.Value) })
        .Aggregate(this.PathPattern, (path, param) => path.Replace(param.Token, param.Value)));

      if (QueryParameters.Count <= 0) return new Uri(this.BaseUri, restPath.ToString());
      restPath.Append(string.IsNullOrEmpty(BaseUri.Query) ? "?" : "&");
      // If parameter value is empty - just add the "name", otherwise "name=value"
      restPath.Append(String.Join("&", QueryParameters.Select(
        x => string.IsNullOrEmpty(x.Value) ?
          Uri.EscapeDataString(x.Key) :
          String.Format("{0}={1}", Uri.EscapeDataString(x.Key), Uri.EscapeDataString(x.Value)))
        .ToArray()));

      return new Uri(this.BaseUri, restPath.ToString());
    }

    public string RelativePath()
    {
      var restPath = new StringBuilder(PathParameters
        .Select(param => new { Token = "{" + param.Key + "}", Value = Uri.EscapeDataString(param.Value) })
        .Aggregate(this.PathPattern, (path, param) => path.Replace(param.Token, param.Value)));
      return restPath.ToString();
    }
    /// <summary>Adds a parameter value.</summary> 
    /// <param name="type">Type of the parameter (must be Path or Query).</param>
    /// <param name="name">Parameter name.</param>
    /// <param name="value">Parameter value.</param>
    public void AddParameter(HttpRequestParameter.ParamType type, string name, string value)
    {
      switch (type)
      {
        case HttpRequestParameter.ParamType.Path:
          if (string.IsNullOrEmpty(value))
          {
            throw new ArgumentException("Path parameters cannot be null or empty.");
          }
          PathParameters.Add(name, value);
          break;
        case HttpRequestParameter.ParamType.Query:
          if (String.IsNullOrEmpty(value)) // don't allow null values on query (empty value is valid)
          {
            break;
          }
          QueryParameters.Add(new KeyValuePair<string, string>(name, value));
          break;
        default:
          throw new ArgumentOutOfRangeException("type");
      }
    }
    /// <summary>
    /// Creates a <seealso cref="System.Net.Http.FormUrlEncodedContent"/> with all the specified parameters in 
    /// input request. It uses reflection to iterate over all properties with
    /// <seealso cref="HttpRequestParameter"/> attribute.
    /// </summary>
    /// <returns>
    /// A <seealso cref="System.Net.Http.FormUrlEncodedContent"/> which contains the all the given object required 
    /// values
    /// </returns>
    public FormUrlEncodedContent CreateFormUrlEncodedContent()
    {
      return new FormUrlEncodedContent(QueryParameters);
    }
  }
}