using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedEntity : FeedSource
  {
    [DataMember(EmitDefaultValue = false)]
    public FeedEntry[] entries { get; set; }
  }
}