using System.Collections.Generic;
using Newtonsoft.Json;

namespace famous.oauth.core
{
  /// <summary>
  /// Calls to Google Api return StandardResponses as Json with
  /// two properties Data, being the return type of the method called
  /// and Error, being any errors that occure.
  /// </summary>
  public sealed class GenericResult<TInnerType>
  {
    /// <summary>May be null if call failed.</summary>
    [JsonProperty("data")]
    public TInnerType Data { get; set; }

    /// <summary>
    /// The error code returned
    /// </summary>
    [JsonProperty("code")]
    public int Code { get; set; }

    /// <summary>
    /// The error message returned
    /// </summary>
    [JsonProperty("reason")]
    public string Reason { get; set; }
  }
}