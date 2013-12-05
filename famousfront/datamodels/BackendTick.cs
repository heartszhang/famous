using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class BackendTick
  {
    [DataMember(EmitDefaultValue = false)]
    public long tick { get; set; }  // nano seconds
    [DataMember(EmitDefaultValue = false)]
    public FeedEntity[] feeds { get; set; }
  }
}