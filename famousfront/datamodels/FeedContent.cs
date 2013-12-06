using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  /*
   type FeedContent struct {
	Uri     string      `json:"uri" bson:"uri"`
	FlowDoc string      `json:"doc,omitempty" bson:"doc,omitempty"`
	Local   string      `json:"local,omitempty" bson:"local,omitempty"`
	Words   uint        `json:"words" bson:"words"`
	Density uint        `json:"density" bson:"density"`
	Links   uint        `json:"links" bson:"links"`
	Status  uint64      `json:"status" bson:"status"`
	Images  []FeedMedia `json:"images" bson:"images"`
	Medias  []FeedMedia `json:"media" bson:"media"`
}
*/
  [DataContract]
  internal class FeedContent
  {
    [DataMember(EmitDefaultValue = false)]
    public string uri
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string doc
    {
      get;
      set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string local
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint words
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint density
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint links
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public ulong status
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] images
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] medias
    {
      get;set;
    }
  }
}