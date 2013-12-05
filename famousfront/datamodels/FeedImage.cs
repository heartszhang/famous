using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedImage  // => ImageCache
  {
    [DataMember(EmitDefaultValue = false)]
    public string mime { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string thumbnail { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string origin { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int code { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int width { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int height { get; set; }
  }
}